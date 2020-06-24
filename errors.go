package httparse

import (
	"errors"
)

//StatusPartial 代表字节数长度不够，需要下次重试，在一部非阻塞系统比较常见
var StatusPartial = errors.New("Partial")
