package users

import (
	"ecoctl/command"
	"errors"
	"github.com/ipaqsa/ecogo"
)

var userOpt ecogo.UserOpt

func initFlags() error {
	requestUserConfigCmd.Flags().StringVarP(&command.GlobalOption.PathToConfig, "output", "o", "", "path to config")
	if err := requestUserConfigCmd.MarkFlagRequired("output"); err != nil {
		return errors.Join(command.ErrMarkRequired, err)
	}

	requestUserConfigCmd.Flags().StringVarP(&userOpt.Name, "user", "u", "", "user name")
	if err := requestUserConfigCmd.MarkFlagFilename("user"); err != nil {
		return errors.Join(command.ErrMarkRequired, err)
	}

	requestUserConfigCmd.Flags().Uint64VarP(&userOpt.SecondsExp, "sec", "s", 0, "expiration time")
	if err := requestUserConfigCmd.MarkFlagFilename("sec"); err != nil {
		return errors.Join(command.ErrMarkRequired, err)
	}

	requestUserConfigCmd.Flags().StringVarP(&userOpt.Role, "role", "", "", "user role")
	if err := requestUserConfigCmd.MarkFlagFilename("role"); err != nil {
		return errors.Join(command.ErrMarkRequired, err)
	}

	requestUserConfigCmd.Flags().StringVarP(&userOpt.Namespace, "namespace", "n", "", "user namespace")
	if err := requestUserConfigCmd.MarkFlagFilename("namespace"); err != nil {
		return errors.Join(command.ErrMarkRequired, err)
	}

	requestAdminConfigCmd.Flags().StringVarP(&command.GlobalOption.AdminTTL, "ttl", "", "", "ttl")
	if err := requestAdminConfigCmd.MarkFlagFilename("ttl"); err != nil {
		return errors.Join(command.ErrMarkRequired, err)
	}
	requestAdminConfigCmd.Flags().StringVarP(&command.GlobalOption.PathToConfig, "output", "o", "", "path to config")
	if err := requestAdminConfigCmd.MarkFlagRequired("output"); err != nil {
		return errors.Join(command.ErrMarkRequired, err)
	}

	return nil
}
