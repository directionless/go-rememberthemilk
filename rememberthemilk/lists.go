package rememberthemilk

import (
	"github.com/pkg/errors"
)

func (rtm *RememberTheMilk) GetList() ([]List, error) {

	resp := &listResponse{}

	if err := rtm.Req("rtm.lists.getList", resp); err != nil {
		return nil, errors.Wrap(err, "failed req")
	}

	return resp.Lists, nil
}
