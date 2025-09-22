package tor

import (
	"os"

	"github.com/cretz/bine/tor"
)

type Instance struct {
	client *tor.Tor
	dialer *tor.Dialer
}

func newInstance() (*Instance, error) {
	if dataDir, err := os.MkdirTemp("/tmp/", "data-dir-"); err != nil {
		return nil, err
	} else {
		t, err := tor.Start(
			nil,
			&tor.StartConf{
				DataDir: dataDir,
				ExtraArgs: []string{
					"--Log", "err file /dev/null",
				},
			},
		)
		if err != nil {
			return nil, err
		}
		dialer, err := t.Dialer(nil, nil)
		if err != nil {
			return nil, err
		}
		return &Instance{
			client: t,
			dialer: dialer,
		}, nil
	}
}
