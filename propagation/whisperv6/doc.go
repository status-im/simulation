// Package whisperv6 implements message propagation simulator based on the go-ethereum (geth) implementation of
// Whisper v6 protocol.
//
// Whisper specification is defined in EIP 627: http://eips.ethereum.org/EIPS/eip-627
//
// Quick link of whisperv6 implementation: https://github.com/ethereum/go-ethereum/tree/master/whisper/whisperv6
//
// This simulator utilizes geth's simulations (https://github.com/ethereum/go-ethereum/tree/master/p2p/simulations) package and starts in-memory whisper services with a given
// network topology.
package whisperv6

//go:generate autoreadme -f
