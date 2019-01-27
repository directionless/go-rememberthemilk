package rememberthemilk

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
)

func (rtm *RememberTheMilk) GetList() error {

	spew.Dump(rtm)
	spew.Dump(rtm.apiToken)

	mv, err := rtm.Req("rtm.lists.getList")
	if err != nil {
		return errors.Wrap(err, "failed req")
	}

	v, err := mv.ValueForPathString("rsp")
	if err != nil {
		return errors.Wrap(err, "Failed to get value")
	}

	spew.Dump(v)

	return nil
}
