package connection

import (
	"net";
	"fmt";
	"os";
	"io";
	"log";
	"redis";
	"protocol";
)

const (
	TCP = "tcp";
	LOCALHOST = "127.0.0.1";
)

// ----------------------------------------------------------------------------
// Connection Spec
// ----------------------------------------------------------------------------

type Spec struct {
	host string;
	port int;
	password string;
	db 	int;
}

func DefaultSpec () *Spec {
	spec := new (Spec);
	//var spec Spec;
	spec.host = LOCALHOST;
	spec.port = 6379;
	return spec;
}
func (spec *Spec) Db(db int) *Spec {
	spec.db = db;
	return spec;
}
func (spec *Spec) Host(host string) *Spec {
	spec.host = host;
	return spec;
}
func (spec *Spec) Port(port int) *Spec {
	spec.port = port;
	return spec;
}
func (spec *Spec) Password(password string) *Spec {
	spec.password = password;
	return spec;
}
func (spec *Spec) Addr () string {
	return fmt.Sprintf("%s:%d", spec.host, spec.port);
}

// ----------------------------------------------------------------------------
// Connection Endpoint
// ----------------------------------------------------------------------------

type Endpoint interface {
	ServiceRequest (cmd *redis.Command, args ...) (protocol.Response, redis.Error);
	Close () os.Error;
}
type _connection struct {
	spec *Spec;
	conn net.Conn;
}
func (hdl _connection) Close() os.Error {
	err := hdl.conn.Close();
	log.Stdout ("Closed connection");
	return err;
}

func (c _connection) ServiceRequest (cmd *redis.Command, args ...) (resp protocol.Response, err redis.Error) {
	
	// TODO: need to consider errors here -- assuming it is always ok to write to Buffer ...
	buff, e1 := protocol.CreateRequestBytes(cmd, args);
	if e1 != nil {
		return nil, redis.NewErrorWithCause(redis.SYSTEM_ERR, "ServiceRequest(): failed to create request buffer", e1);
	}
	
	e2 := c.sendRequest(c.conn, buff);
	if e2 != nil {
		return nil, redis.NewErrorWithCause(redis.SYSTEM_ERR, "ServiceRequest(): failed to send request", e2);
	}
	
	resp, e3 := protocol.GetResponse(c.conn, cmd);
	if e3 != nil {
		return nil, redis.NewErrorWithCause(redis.SYSTEM_ERR, "ServiceRequest(): failed to get response", e3);
	}
	
	if resp.IsError() {
		log.Stderr("REDIS ERROR: ", resp.GetMessage());
		return nil, redis.NewRedisError(resp.GetMessage());
	}
	return;
}

func OpenNew (spec *Spec) (c Endpoint, err os.Error) {
	hdl := new(_connection);
	addr := spec.Addr();
	hdl.conn, err = net.Dial(TCP, "", addr);	
	switch {
		case err != nil:
			err = redis.NewErrorWithCause(redis.SYSTEM_ERR, "Could not open connection", err);
		default:
			log.Stdout("Opened connection to ", addr);
			c = hdl;
	}
	return;
}

// ----------------------------------------------------------------------------
// internal ops
// ----------------------------------------------------------------------------

func (hdl _connection) sendRequest (conn io.Writer, data []byte) os.Error {
	n, e1 := conn.Write(data);
	if e1 != nil {
		log.Stderr ("error on Write: ", e1);
	}
	if n < 6 {
		log.Stderr ("didn't write the whole data: ", n);
	}
	else {
//		log.Stderr ("wrote data: ", n);
	}
	return e1;
}

