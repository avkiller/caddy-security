// Copyright 2022 Paul Greenberg greenpau@outlook.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package security

import (
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/greenpau/aaasf"
	"github.com/greenpau/aaasf/pkg/authz"
	"github.com/greenpau/aaasf/pkg/errors"
	"github.com/greenpau/caddy-security/pkg/util"
)

const (
	authzPrefix = "security.authorization"
)

// parseCaddyfileAuthorization parses authorization configuration.
//
// Syntax:
//
//   authorization portal <name> {
//   }
//
func parseCaddyfileAuthorization(d *caddyfile.Dispenser, repl *caddy.Replacer, cfg *aaasf.Config) error {
	var rootDirective string
	args := util.FindReplaceAll(repl, d.RemainingArgs())
	if len(args) != 2 {
		return d.ArgErr()
	}
	switch args[0] {
	case "policy":
		p := &authz.PolicyConfig{Name: args[1]}
		for nesting := d.Nesting(); d.NextBlock(nesting); {
			k := d.Val()
			v := util.FindReplaceAll(repl, d.RemainingArgs())
			rootDirective = mkcp(authzPrefix, args[0], k)
			switch k {
			case "crypto":
				if err := parseCaddyfileAuthorizationCrypto(d, repl, p, rootDirective, v); err != nil {
					return err
				}
			case "acl":
				if err := parseCaddyfileAuthorizationACL(d, repl, p, rootDirective, v); err != nil {
					return err
				}
			case "allow", "deny":
				if err := parseCaddyfileAuthorizationACLShortcuts(d, repl, p, rootDirective, k, v); err != nil {
					return err
				}
			case "bypass":
				if err := parseCaddyfileAuthorizationBypass(d, repl, p, rootDirective, v); err != nil {
					return err
				}
			case "enable", "disable", "validate", "set", "with":
				if err := parseCaddyfileAuthorizationMisc(d, repl, p, rootDirective, k, v); err != nil {
					return err
				}
			case "inject":
				if err := parseCaddyfileAuthorizationHeaderInjection(d, repl, p, rootDirective, v); err != nil {
					return err
				}
			default:
				return errors.ErrMalformedDirective.WithArgs(rootDirective, v)
			}
		}
		if err := cfg.AddAuthorizationPolicy(p); err != nil {
			return err
		}
	default:
		return errors.ErrMalformedDirective.WithArgs(authzPrefix, args)
	}
	return nil
}
