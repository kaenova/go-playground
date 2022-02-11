package list

import "golang.org/x/net/context"

type Server struct {
}

func (s *Server) TotalElement(ctx context.Context, req *ReqBuffer) (*ResBuffer, error) {
	var res ResBuffer

	res.TotalElement = int32(len(req.Text))

	return &res, nil
}
