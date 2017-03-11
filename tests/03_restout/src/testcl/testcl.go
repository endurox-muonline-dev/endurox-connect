package main

import (
	"errors"
	"fmt"
	"os"
	u "ubftab"

	atmi "github.com/endurox-dev/endurox-go"
)

/*
#include <signal.h>
*/
import "C"

const (
	ProgSection = "testcl"
)

var MSomeConfigFlag string = ""
var MSomeOtherConfigFlag int = 0

func juerrors(ac *atmi.ATMICtx) error {

	for i := 0; i < 10000; i++ {

		buf, err := ac.NewUBF(1024)

		if err != nil {
			return errors.New(err.Error())
		}

		//Set some field for call
		buf.BChg(u.T_CHAR_FLD, 0, "A")
		buf.BChg(u.T_SHORT_FLD, 0, "123")
		buf.BChg(u.T_LONG_FLD, 0, i)
		buf.BChg(u.T_FLOAT_FLD, 0, "1.33")
		buf.BChg(u.T_DOUBLE_FLD, 0, "4444.3333")
		buf.BChg(u.T_STRING_FLD, 0, "HELLO")
		buf.BChg(u.T_CARRAY_FLD, 0, "WORLD")

		//Call the server
		if _, err := ac.TpCall("JUERRORS", buf, 0); nil != err {
			return errors.New(err.Error())
		}

		//Test the response...
		stmt := fmt.Sprintf("T_SHORT_2_FLD==123 "+
			"&& T_LONG_2_FLD==%d"+
			"&& T_CHAR_2_FLD=='A'"+
			"&& T_FLOAT_2_FLD==1.33"+
			"&& T_DOUBLE_2_FLD==4444.3333"+
			"&& T_STRING_2_FLD=='HELLO'"+
			"&& T_CARRAY_2_FLD=='WORLD'", i)

		if res, err := buf.BQBoolEv(stmt); !res || nil != err {
			if nil != err {
				return errors.New(fmt.Sprintf("juerrors: Expression "+
					"failed: %s", err.Error()))
			} else {
				return errors.New("juerrors: Expression is FALSE!: %s")
			}
		}
	}

	return nil
}

//Run the listener
func apprun(ac *atmi.ATMICtx) error {

	//Do some work here

	if len(os.Args) != 2 {
		return errors.New(fmt.Sprintf("usage: %s <command>",
			os.Args[0]))
	}

	cmd := os.Args[1]

	ac.TpLogInfo("Got command: [%s]", cmd)

	//These are projection on 01_restin/runtime/conf/restin.ini cases
	switch cmd {
	case "juerrors":
		return juerrors(ac)
	default:
		return errors.New(fmt.Sprintf("Invalid test case: [%s]", cmd))
	}

}

//Init function
//@param ac	ATMI context
//@return error (if erro) or nil
func appinit(ac *atmi.ATMICtx) error {

	if err := ac.TpInit(); err != nil {
		return errors.New(err.Error())
	}

	return nil
}

//Un-init & Terminate the application
//@param ac	ATMI Context
//@param restCode	Return code. atmi.FAIL (-1) or atmi.SUCCEED(0)
func unInit(ac *atmi.ATMICtx, retCode int) {

	ac.TpTerm()
	ac.FreeATMICtx()
	os.Exit(retCode)
}

//Cliet process main entry
func main() {

	ac, errA := atmi.NewATMICtx()

	if nil != errA {
		fmt.Fprintf(os.Stderr, "Failed to allocate cotnext %d:%s!\n",
			errA.Code(), errA.Message())
		os.Exit(atmi.FAIL)
	}

	if err := appinit(ac); nil != err {
		ac.TpLogError("Failed to init: %s", err)
		os.Exit(atmi.FAIL)
	}

	ac.TpLogWarn("Init complete, processing...")

	if err := apprun(ac); nil != err {
		ac.TpLogError("Got error: [%s]", err.Error())
		unInit(ac, atmi.FAIL)
	}

	unInit(ac, atmi.SUCCEED)
}