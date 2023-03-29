package executer

import "context"

func Run(ctx context.Context, path string, args ...string) error {
	binaries, err := prepare(ctx, path)
	if err != nil {
		return err
	}

	if err := executeAction(ctx, binaries, args...); err != nil {
		return err
	}

	return nil
}
