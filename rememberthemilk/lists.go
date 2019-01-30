package rememberthemilk

import (
	"github.com/davecgh/go-spew/spew"
	"github.com/pkg/errors"
)

func (rtm *RememberTheMilk) GetList() error {

	lists := &ListResponse{}

	//spew.Dump(rtm)
	//spew.Dump(rtm.apiToken)

	if err := rtm.Req("rtm.lists.getList", lists); err != nil {
		return errors.Wrap(err, "failed req")
	}

	spew.Dump(lists)

	return nil
}
