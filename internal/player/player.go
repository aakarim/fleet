package player

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/aakarim/fleet/internal/host"
	"github.com/bramvdbogaerde/go-scp"
)

type Player struct {
	starshPath string
}

func NewPlayer(starshPath string) *Player {
	return &Player{
		starshPath: starshPath,
	}
}

func (p *Player) Play(ctx context.Context, host *host.Host, playbook io.Reader) error {
	if _, err := host.Connect(); err != nil {
		return fmt.Errorf("error connecting to host(%s): %w", host.Addr, err)
	}

	defer host.Close()

	fmt.Println("Connected to host")

	if err := p.updateStarsh(ctx, host); err != nil {
		return fmt.Errorf("error updating Starshell on host: %w", err)
	}

	if err := p.uploadPlaybook(ctx, host, playbook); err != nil {
		return fmt.Errorf("error uploading playbook: %w", err)
	}

	if err := p.runPlaybook(ctx, host); err != nil {
		return fmt.Errorf("error running playbook: %w", err)
	}
	return nil
}

func (p *Player) runPlaybook(ctx context.Context, host *host.Host) error {
	session, err := host.Connection.NewSession()
	if err != nil {
		return fmt.Errorf("error creating session: %w", err)
	}

	session.Stdout = os.Stdout
	session.Stderr = os.Stderr

	if err := session.Run("/tmp/starsh /tmp/playbook.star"); err != nil {
		return fmt.Errorf("error running playbook: %w", err)
	}

	return nil
}

func (p *Player) updateStarsh(ctx context.Context, host *host.Host) error {
	// copy over starshell
	// TODO: gather facts, check existing starshell and only upload if necessary
	client, err := scp.NewClientBySSH(host.Connection)
	if err != nil {
		return fmt.Errorf("NewClientBySSH(): %w", err)
	}

	f, err := os.Open(p.starshPath)
	if err != nil {
		return fmt.Errorf("error opening Starshell: %w", err)
	}

	defer f.Close()

	if err := client.Connect(); err != nil {
		return fmt.Errorf("error connecting to host(%s): %w", host.Addr, err)
	}

	defer client.Close()

	if err := client.CopyFromFile(ctx, *f, "/tmp/starsh", "0755"); err != nil {
		return fmt.Errorf("error copying file: %w", err)
	}

	return nil
}

func (p *Player) uploadPlaybook(ctx context.Context, host *host.Host, playbook io.Reader) error {
	// copy over playbook
	client, err := scp.NewClientBySSH(host.Connection)
	if err != nil {
		return err
	}

	if err := client.Connect(); err != nil {
		return err
	}

	defer client.Close()

	pb, err := io.ReadAll(playbook)
	if err != nil {
		return err
	}

	if err := client.Copy(ctx, bytes.NewReader(pb), "/tmp/playbook.star", "0644", int64(len(pb))); err != nil {
		return err
	}

	return nil
}
