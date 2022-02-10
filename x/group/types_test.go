package group_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/x/group"

	"github.com/cosmos/cosmos-sdk/x/group/internal/math"

	"github.com/stretchr/testify/require"
)

func TestPercentageDecisionPolicyAllow(t *testing.T) {
	testCases := []struct {
		name           string
		policy         *group.PercentageDecisionPolicy
		tally          *group.Tally
		totalPower     string
		votingDuration time.Duration
		result         group.DecisionPolicyResult
	}{
		{
			"YesCount percentage > decision policy percentage",
			&group.PercentageDecisionPolicy{
				Percentage: "0.5",
				Timeout:    time.Second * 100,
			},
			&group.Tally{
				YesCount:     "2",
				NoCount:      "0",
				AbstainCount: "0",
				VetoCount:    "0",
			},
			"3",
			time.Duration(time.Second * 50),
			group.DecisionPolicyResult{
				Allow: true,
				Final: true,
			},
		},
		{
			"YesCount percentage == decision policy percentage",
			&group.PercentageDecisionPolicy{
				Percentage: "0.5",
				Timeout:    time.Second * 100,
			},
			&group.Tally{
				YesCount:     "2",
				NoCount:      "0",
				AbstainCount: "0",
				VetoCount:    "0",
			},
			"4",
			time.Duration(time.Second * 50),
			group.DecisionPolicyResult{
				Allow: true,
				Final: true,
			},
		},
		{
			"YesCount percentage < decision policy percentage",
			&group.PercentageDecisionPolicy{
				Percentage: "0.5",
				Timeout:    time.Second * 100,
			},
			&group.Tally{
				YesCount:     "1",
				NoCount:      "0",
				AbstainCount: "0",
				VetoCount:    "0",
			},
			"3",
			time.Duration(time.Second * 50),
			group.DecisionPolicyResult{
				Allow: false,
				Final: false,
			},
		},
		{
			"sum percentage (YesCount + undecided votes percentage) < decision policy percentage",
			&group.PercentageDecisionPolicy{
				Percentage: "0.5",
				Timeout:    time.Second * 100,
			},
			&group.Tally{
				YesCount:     "1",
				NoCount:      "2",
				AbstainCount: "0",
				VetoCount:    "0",
			},
			"3",
			time.Duration(time.Second * 50),
			group.DecisionPolicyResult{
				Allow: false,
				Final: true,
			},
		},
		{
			"sum percentage = decision policy percentage",
			&group.PercentageDecisionPolicy{
				Percentage: "0.5",
				Timeout:    time.Second * 100,
			},
			&group.Tally{
				YesCount:     "1",
				NoCount:      "2",
				AbstainCount: "0",
				VetoCount:    "0",
			},
			"4",
			time.Duration(time.Second * 50),
			group.DecisionPolicyResult{
				Allow: false,
				Final: false,
			},
		},
		{
			"sum percentage > decision policy percentage",
			&group.PercentageDecisionPolicy{
				Percentage: "0.5",
				Timeout:    time.Second * 100,
			},
			&group.Tally{
				YesCount:     "1",
				NoCount:      "0",
				AbstainCount: "0",
				VetoCount:    "0",
			},
			"3",
			time.Duration(time.Second * 50),
			group.DecisionPolicyResult{
				Allow: false,
				Final: false,
			},
		},
		{
			"decision policy timeout <= voting duration",
			&group.PercentageDecisionPolicy{
				Percentage: "0.5",
				Timeout:    time.Second * 10,
			},
			&group.Tally{
				YesCount:     "2",
				NoCount:      "0",
				AbstainCount: "0",
				VetoCount:    "0",
			},
			"3",
			time.Duration(time.Second * 50),
			group.DecisionPolicyResult{
				Allow: false,
				Final: true,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			percentage, err := math.NewPositiveDecFromString(tc.policy.Percentage)
			require.NoError(t, err)
			yesCount, err := math.NewNonNegativeDecFromString(tc.tally.YesCount)
			require.NoError(t, err)
			totalPowerDec, err := math.NewNonNegativeDecFromString(tc.totalPower)
			require.NoError(t, err)
			totalCounts, err := tc.tally.TotalCounts()
			require.NoError(t, err)
			undecided, err := math.SubNonNegative(totalPowerDec, totalCounts)
			require.NoError(t, err)
			sum, err := yesCount.Add(undecided)
			require.NoError(t, err)
			yesPercentage, err := yesCount.Quo(totalPowerDec)
			require.NoError(t, err)
			sumPercentage, err := sum.Quo(totalPowerDec)
			require.NoError(t, err)
			fmt.Println("------------")
			fmt.Println(sumPercentage, percentage)
			fmt.Println(yesPercentage, percentage)
			// panic("")
			policyResult, err := tc.policy.Allow(*tc.tally, tc.totalPower, tc.votingDuration)
			require.NoError(t, err)
			require.Equal(t, tc.result, policyResult)
		})
	}
}

func TestThresholdDecisionPolicyAllow(t *testing.T) {
	testCases := []struct {
		name           string
		policy         *group.ThresholdDecisionPolicy
		tally          *group.Tally
		totalPower     string
		votingDuration time.Duration
		result         group.DecisionPolicyResult
	}{
		{
			"YesCount >= threshold decision policy",
			&group.ThresholdDecisionPolicy{
				Threshold: "3",
				Timeout:   time.Second * 100,
			},
			&group.Tally{
				YesCount:     "3",
				NoCount:      "0",
				AbstainCount: "0",
				VetoCount:    "0",
			},
			"3",
			time.Duration(time.Second * 50),
			group.DecisionPolicyResult{
				Allow: true,
				Final: true,
			},
		},
		{
			"YesCount < threshold decision policy",
			&group.ThresholdDecisionPolicy{
				Threshold: "3",
				Timeout:   time.Second * 100,
			},
			&group.Tally{
				YesCount:     "1",
				NoCount:      "0",
				AbstainCount: "0",
				VetoCount:    "0",
			},
			"3",
			time.Duration(time.Second * 50),
			group.DecisionPolicyResult{
				Allow: false,
				Final: false,
			},
		},
		{
			"sum votes < threshold decision policy",
			&group.ThresholdDecisionPolicy{
				Threshold: "3",
				Timeout:   time.Second * 100,
			},
			&group.Tally{
				YesCount:     "1",
				NoCount:      "0",
				AbstainCount: "0",
				VetoCount:    "0",
			},
			"2",
			time.Duration(time.Second * 50),
			group.DecisionPolicyResult{
				Allow: false,
				Final: true,
			},
		},
		{
			"sum votes >= threshold decision policy",
			&group.ThresholdDecisionPolicy{
				Threshold: "3",
				Timeout:   time.Second * 100,
			},
			&group.Tally{
				YesCount:     "1",
				NoCount:      "0",
				AbstainCount: "0",
				VetoCount:    "0",
			},
			"3",
			time.Duration(time.Second * 50),
			group.DecisionPolicyResult{
				Allow: false,
				Final: false,
			},
		},
		{
			"decision policy timeout <= voting duration",
			&group.ThresholdDecisionPolicy{
				Threshold: "3",
				Timeout:   time.Second * 10,
			},
			&group.Tally{
				YesCount:     "3",
				NoCount:      "0",
				AbstainCount: "0",
				VetoCount:    "0",
			},
			"3",
			time.Duration(time.Second * 50),
			group.DecisionPolicyResult{
				Allow: false,
				Final: true,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			policyResult, err := tc.policy.Allow(*tc.tally, tc.totalPower, tc.votingDuration)
			require.NoError(t, err)
			require.Equal(t, tc.result, policyResult)
		})
	}
}