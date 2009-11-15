package protocol

import (
	"os";
	"io";
	"bufio";
	"strconv";
	"bytes";
	"strings";
	"log";
	"reflect";
	"fmt";
	"redis";
 	ARGS "reflect_utils";
)
// ----------------------------------------------------------------------------
// Wire
// ----------------------------------------------------------------------------

// protocol's special bytes
const 
(
	CR_BYTE 	byte= byte('\r');
	LF_BYTE			= byte('\n');
	SPACE_BYTE 		= byte(' ');
	ERR_BYTE 		= byte(45);
	OK_BYTE 		= byte(43);
	COUNT_BYTE 		= byte(42);
	SIZE_BYTE 		= byte(36);
	NUM_BYTE 		= byte(58);
	FALSE_BYTE 		= byte(48);
	TRUE_BYTE 		= byte(49);
)

type CtlBytes []byte;
var CRLF 		CtlBytes = make([]byte, 2);
var WHITESPACE 	CtlBytes = make([]byte, 1);

func init () {
	CRLF[0] = CR_BYTE; CRLF[1] = LF_BYTE;
	WHITESPACE [0] = SPACE_BYTE;
}

// ----------------------------------------------------------------------------
// Services
// ----------------------------------------------------------------------------

// TODO: tedious but need to check for errors on all buffer writes ..
//
func CreateRequestBytes (cmd *redis.Command, v ...) ([]byte, os.Error) {

	args := reflect.NewValue(v).(*reflect.StructValue);
	cmd_bytes := strings.Bytes(cmd.Code);
	
	buffer := bytes.NewBuffer(cmd_bytes);

	switch cmd.ReqType {

	case redis.NO_ARG:
	
	case redis.KEY:
		buffer.Write(WHITESPACE);	
		key, _ := ARGS.GetByteArrayAtIndex (args, 0);
		log.Stdout("redis.Key => ", key);
		buffer.Write(key);
		
	case 
		redis.KEY_KEY, 
		redis.KEY_NUM, 
		redis.KEY_SPEC:
		
		buffer.Write(WHITESPACE);	
		key, _ := ARGS.GetByteArrayAtIndex (args, 0);
		buffer.Write(key);
		buffer.Write(WHITESPACE);	
		key2, _ := ARGS.GetByteArrayAtIndex (args, 1);
		buffer.Write(key2);
	
	case redis.KEY_NUM_NUM:
	
		buffer.Write(WHITESPACE);	
		key, _ := ARGS.GetByteArrayAtIndex (args, 0);
		buffer.Write(key);
		buffer.Write(WHITESPACE);	
		num1, _ := ARGS.GetByteArrayAtIndex (args, 1);
		buffer.Write(num1);
		buffer.Write(WHITESPACE);	
		num2, _ := ARGS.GetByteArrayAtIndex (args, 2);
		buffer.Write(num2);
	
	case redis.KEY_VALUE:
		buffer.Write(WHITESPACE);	
		key, _ := ARGS.GetByteArrayAtIndex (args, 0);
		buffer.Write(key);
		buffer.Write(WHITESPACE);
		value, _ := ARGS.GetByteArrayAtIndex (args, 1);;
		len := fmt.Sprintf("%d", len(value));
		buffer.Write(strings.Bytes(len)); 
		buffer.Write(CRLF);
		buffer.Write(value);
	
	case 
		redis.KEY_IDX_VALUE,
		redis.KEY_KEY_VALUE:
		
	case redis.KEY_CNT_VALUE:
	
	case redis.MULTI_KEY:
	}
	
	buffer.Write(CRLF);	
	
//	log.Stdout ("** \n", buffer);

	return buffer.Bytes(), nil;
}

// Gets the response to the command.
// Any errors (whether runtime or bugs) are returned as os.Error.
// The returned response (regardless of flavor) may have (application level)
// errors as sent from Redis server.
//
func GetResponse (reader io.Reader, cmd *redis.Command) (resp Response, err os.Error) {
	switch cmd.RespType {
	case redis.BOOLEAN:
	    resp, err = getBooleanResponse(reader, cmd);
	case redis.BULK:
	    resp, err = getBulkResponse (reader, cmd);
	case redis.MULTI_BULK:
	    resp, err = getMultiBulkResponse (reader, cmd);
	case redis.NUMBER:
	    resp, err = getNumberResponse (reader, cmd);
	case redis.STATUS:
	    resp, err = getStatusResponse (reader, cmd);
	case redis.STRING:
	    resp, err = getStringResponse (reader, cmd);
//	case redis.VIRTUAL:
//	    resp, err = getVirtualResponse ();
	}
	return;
}

// ----------------------------------------------------------------------------
// internal ops
// ----------------------------------------------------------------------------

func getStatusResponse (conn io.Reader, cmd *redis.Command) (resp Response, e os.Error) {
	buff, error, fault := readLine(conn);
	if fault == nil {
		line := bytes.NewBuffer(buff).String();
		resp = newStatusResponse(line, error);
	}
	return resp, fault;
}

func getBooleanResponse (conn io.Reader, cmd *redis.Command) (resp Response, e os.Error) {
	buff, error, fault := readLine(conn);
	if fault == nil {
		if !error {
			b := buff[1] == TRUE_BYTE;
			resp = newBooleanResponse(b, error);
		}
		else { resp = newStatusResponse(bytes.NewBuffer(buff).String(), error); }
	}
	return resp, fault;
}

func getStringResponse (conn io.Reader, cmd *redis.Command) (resp Response, e os.Error) {
	buff, error, fault := readLine(conn);
	if fault == nil {
		if !error {
			buff = buff[1: len(buff)];
			str := bytes.NewBuffer(buff).String();
			resp = newStringResponse(str, error);
		}
		else { resp = newStatusResponse(bytes.NewBuffer(buff).String(), error); }
	}
	return resp, fault;
}
func getNumberResponse (conn io.Reader, cmd *redis.Command)  (resp Response, e os.Error) {

	buff, error, fault := readLine(conn);
	if fault == nil {
		if !error {
			buff = buff[1: len(buff)];
			numrep := bytes.NewBuffer(buff).String();
			num, err := strconv.Atoi64(numrep);
			if err == nil {  resp = newNumberResponse(num, error);  }
			else { e = os.NewError("<BUG> Expecting a int64 number representation here: " + err.String()); }
		}
		else { resp = newStatusResponse(bytes.NewBuffer(buff).String(), error); }
	}
	return resp, fault;
}

func getBulkResponse (conn io.Reader, cmd *redis.Command) (Response, os.Error) {
/*
	num, error, fault := readCtlNum(conn, true, SIZE_BYTE);
	if fault == nil {
		if !error {
			buff = buff[1: len(buff)];
			numrep := bytes.NewBuffer(buff).String();
			num, err := strconv.Atoi64(numrep);
			if err == nil {  resp = newNumberResponse(num, error);  }
			else { e = os.NewError("<BUG> Expecting a int64 number representation here: " + err.String()); }
		}
		else { resp = newStatusResponse(bytes.NewBuffer(buff).String(), error); }
	}
*/	
	num, _, _ := readCtlNum(conn, true, SIZE_BYTE);
	log.Stderr("bulk data size: ", num);
	return nil, os.NewError("<BUG> Not implemented");;
}
func getMultiBulkResponse (conn io.Reader, cmd *redis.Command) (Response, os.Error) {
	return nil, os.NewError("GetMultiBullkResponse NOT IMPLEMENTED!");
}


// ----------------------------------------------------------------------------
// Response
// ----------------------------------------------------------------------------

type Response interface {
	IsError () bool;
	GetMessage() string;
	GetBooleanValue () bool;
	GetNumberValue() int64;
	GetStringValue () string;
	GetBulkData() []byte;
	GetMultiBulkData() [][]byte;
}
type _response struct {
	isError 		bool;
	msg     		string;
	boolval 		bool;
	numval  		int64;
	stringval 		string;
	bulkdata 		[]byte;
	multibulkdata 	[][]byte;
}
func (r *_response) IsError () bool { return r.isError;}
func (r *_response) GetMessage() string {return r.msg;}
func (r *_response) GetBooleanValue () bool {return r.boolval;}
func (r *_response) GetNumberValue() int64 {return r.numval;}
func (r *_response) GetStringValue () string {return r.stringval;}
func (r *_response) GetBulkData() []byte {return r.bulkdata;}
func (r *_response) GetMultiBulkData() [][]byte {return r.multibulkdata;}
func newAndInitResponse(isError bool) (r *_response) {
	r = new(_response);
	r.isError = isError;
	r.bulkdata = nil;
	r.multibulkdata = nil;
	return;
}
func newStatusResponse (msg string, isError bool) Response {
	r := newAndInitResponse(isError);
	r.msg = msg;
	return r;
}
func newBooleanResponse (val bool, isError bool) Response {
	r := newAndInitResponse(isError);
	r.boolval = val;
	return r;
}
func newNumberResponse (val int64, isError bool) Response {
	r := newAndInitResponse(isError);
	r.numval = val;
	return r;
}
func newStringResponse (val string, isError bool) Response {
	r := newAndInitResponse(isError);
	r.stringval = val;
	return r;
}
func newBulkResponse (val []byte, isError bool) Response {
	r := newAndInitResponse(isError);
	r.bulkdata = val;
	return r;
}
func newMultiBulkResponse (val [][]byte, isError bool) Response {
	r := newAndInitResponse(isError);
	r.multibulkdata = val;
	return r;
}

// ----------------------------------------------------------------------------
// Protocol i/o
// ----------------------------------------------------------------------------

// reads all bytes upto CR-LF.  (Will eat those last two bytes)
// return the line []byte up to CR-LF
// error returned is NOT ("-ERR ...").  If there is a Redis error
// that is in the line buffer returned

func readToCRLF (conn io.Reader) (buffer []byte, err os.Error) {
	reader := bufio.NewReader(conn);
	var buf []byte;
	buf, err = reader.ReadBytes(CR_BYTE);
	if err == nil {
		var b byte;
		b, err = reader.ReadByte();
		if err != nil { return; }
		if b != LF_BYTE { 
			err = os.NewError("<BUG> Expecting a Linefeed byte here!");
		}
		log.Stderr("readToCRLF: ", buf);
		buffer = buf[0 : len(buf) - 1];
	}
	return;
}

func readLine (conn io.Reader) (buf []byte, error bool, fault os.Error) {
	buf, fault = readToCRLF (conn);
	if fault == nil {
		error = buf[0] == ERR_BYTE;
	}
	return;
}

func readCtlNum (conn io.Reader, checkForError bool, ctlByte byte) (ctlNum int64, error redis.Error, fault os.Error) {
	buff, isError, fault := readLine(conn);
	if fault == nil {
		if isError && checkForError {
			log.Stderr("we have a redis error in readCtlNum");
			error = redis.NewRedisError(bytes.NewBuffer(buff).String());
		}
		else {
			if buff[0] == ctlByte {
				buff = buff[1: len(buff)];
				numrep := bytes.NewBuffer(buff).String();
				ctlNum, fault = strconv.Atoi64(numrep);
			}
			else {
				log.Stderr("we have a fault in readCtlNum");
				fault = os.NewError ("<BUG> Expecting a ctl byte here!"); 
			}
		}
	}
	return;
}
