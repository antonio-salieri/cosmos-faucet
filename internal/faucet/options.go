package faucet

import (
	"strings"
)

const (
	DefaultAppCli        = "liquidityd"
	DefaultKeyName       = "faucet"
	DefaultDenoms        = "atom,stake,scrt,iris,band,kava,usdt,luna"
	DefaultCreditAmount  = 300
	DefaultMaximumCredit = 300
)

func defaultOptions() *Options {
	return &Options{
		AppCli:       DefaultAppCli,
		KeyName:      DefaultKeyName,
		Denoms:       strings.Split(DefaultDenoms, ","),
		CreditAmount: DefaultCreditAmount,
		MaxCredit:    DefaultMaximumCredit,
	}
}

type Options struct {
	AppCli          string
	KeyringPassword string
	KeyName         string
	KeyMnemonic     string
	Denoms          []string
	CreditAmount    uint64
	MaxCredit       uint64
}

type Option func(*Options)

func CliName(s string) Option {
	return func(opts *Options) {
		opts.AppCli = s
	}
}

func KeyringPassword(s string) Option {
	return func(opts *Options) {
		opts.KeyringPassword = s
	}
}

func KeyName(s string) Option {
	return func(opts *Options) {
		opts.KeyName = s
	}
}

func WithMnemonic(s string) Option {
	return func(opts *Options) {
		opts.KeyMnemonic = s
	}
}

func Denoms(s string) Option {
	return func(opts *Options) {
		opts.Denoms = strings.Split(s, ",")
	}
}

func CreditAmount(v uint64) Option {
	return func(opts *Options) {
		opts.CreditAmount = v
	}
}

func MaxCredit(v uint64) Option {
	return func(opts *Options) {
		opts.MaxCredit = v
	}
}
