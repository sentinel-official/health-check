package record

import (
	"github.com/gin-gonic/gin"
)

type RequestGetRecords struct {
}

func NewRequestGetRecords(c *gin.Context) (req *RequestGetRecords, err error) {
	req = &RequestGetRecords{}

	return req, nil
}

type RequestGetRecord struct {
	URI struct {
		Addr string `uri:"addr"`
	}
}

func NewRequestGetRecord(c *gin.Context) (req *RequestGetRecord, err error) {
	req = &RequestGetRecord{}
	if err = c.ShouldBindUri(&req.URI); err != nil {
		return nil, err
	}

	return req, nil
}
