package msg

import (
	"io"

	"github.com/elastos/Elastos.ELA.Utility/common"
)

type DPosProposalVoteSlot struct {
	Hash  common.Uint256
	Votes []DPosProposalVote
}

func (p *DPosProposalVoteSlot) TryAppend(v DPosProposalVote) bool {
	if p.Hash.IsEqual(v.Proposal.BlockHash) {
		p.Votes = append(p.Votes, v)
		return true
	}
	return false
}

func (p *DPosProposalVoteSlot) Serialize(w io.Writer) error {
	if err := p.Hash.Serialize(w); err != nil {
		return err
	}

	if err := common.WriteUint64(w, uint64(len(p.Votes))); err != nil {
		return err
	}

	for _, sign := range p.Votes {
		if err := sign.Serialize(w); err != nil {
			return err
		}
	}

	return nil
}

func (p *DPosProposalVoteSlot) Deserialize(r io.Reader) error {
	if err := p.Hash.Deserialize(r); err != nil {
		return err
	}

	signCount, err := common.ReadUint64(r)
	if err != nil {
		return err
	}
	p.Votes = make([]DPosProposalVote, signCount)

	for i := uint64(0); i < signCount; i++ {
		err := p.Votes[i].Deserialize(r)
		if err != nil {
			return err
		}
	}

	return nil
}