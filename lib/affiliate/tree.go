package affiliate

import (
	"github.com/jinzhu/copier"
	"github.com/sirupsen/logrus"
	"gitlab.com/snap-clickstaff/go-common/comlogging"
	"gitlab.com/snap-clickstaff/go-common/comutils"
	"gitlab.com/snap-clickstaff/torque-go/database"
	"gitlab.com/snap-clickstaff/torque-go/database/models"
	"gitlab.com/snap-clickstaff/torque-go/lib/constants"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
	"gitlab.com/snap-clickstaff/torque-go/lib/trading/tradingbalance"
)

type TreeUser struct {
	ID             meta.UID `json:"id"`
	ParentID       meta.UID `json:"parent_id"`
	Code           string   `json:"code"`
	Username       string   `json:"username"`
	FirstName      string   `json:"first_name"`
	LastName       string   `json:"last_name"`
	Email          string   `json:"email"`
	ReferralCode   string   `json:"referral_code"`
	TierType       int      `json:"tier_type"`
	TierTypeStatus int      `json:"tier_type_status"`

	CreateTime int64 `json:"create_time"`
}

type TreeNode struct {
	User           TreeUser                      `json:"user"`
	CoinBalanceMap tradingbalance.CoinBalanceMap `json:"coin_balance_map"`
	tier           *TreeNodeTier
	stats          *TreeNodeStats

	Parent   *TreeNode   `json:"-"`
	Children []*TreeNode `json:"children"`
}

type TreeNodeMap map[meta.UID]*TreeNode

type Tree struct {
	root    *TreeNode
	nodeMap TreeNodeMap
}

func (tn *TreeNode) GetTier() TreeNodeTier {
	if tn.tier == nil {
		tier := tn.findTier()
		tn.tier = &tier
	}

	return *tn.tier
}

func (tn *TreeNode) GetTierCode() string {
	return tn.GetTier().Meta.Code
}

func (tn *TreeNode) findTier() TreeNodeTier {
	if tn.User.TierType == constants.TierMetaGovernor.ID {
		return TreeNodeTier{
			Meta: constants.TierMetaGovernor,
		}

	} else if tn.User.TierTypeStatus == constants.TierTypeStatusFixed {
		if fixedTierMeta, ok := constants.TierMetaMap[tn.User.TierType]; ok {
			return TreeNodeTier{
				Meta:    fixedTierMeta,
				IsFixed: true,
			}
		}
	}

	if len(tn.Children) < TierPromotionLowerCount {
		return TreeNodeTier{
			Meta: constants.TierMetaAgent,
		}
	}

	logger := comlogging.GetLogger()
	for i := len(constants.TierMetaLeaderOrderedList) - 1; i >= 0; i-- {
		tierMeta := constants.TierMetaLeaderOrderedList[i]

		var tierNodes []*TreeNode
		for _, child := range tn.Children {
			tierNode := child.GetStats().GetNodeByTier(tierMeta.Code)
			if tierNode == nil {
				continue
			}

			tierNodes = append(tierNodes, tierNode)
		}

		if len(tierNodes) >= TierPromotionLowerCount {
			upperTierMeta, err := GetUpperTierMeta(tierMeta.Code)
			if err != nil {
				logger.Errorf("tier `%s` doesn't has an upper tier", tierMeta.Code)
				continue
			}

			return TreeNodeTier{
				Meta:      upperTierMeta,
				BaseNodes: tierNodes,
			}
		}
	}

	var marketLeaderNodes []*TreeNode
	for _, child := range tn.Children {
		if child.GetStats().GetMaxChainLength() >= TierPromotionMarketLeaderDepth {
			marketLeaderNodes = append(marketLeaderNodes, child)
		}
	}
	if len(marketLeaderNodes) >= TierPromotionLowerCount {
		return TreeNodeTier{
			Meta:      constants.TierMetaSeniorPartner,
			BaseNodes: marketLeaderNodes,
		}
	}

	return TreeNodeTier{
		Meta: constants.TierMetaAgent,
	}
}

func (tn *TreeNode) GetStats() *TreeNodeStats {
	if tn.stats == nil {
		stats := tn.genStats()
		tn.stats = &stats
	}

	return tn.stats
}

func (tn *TreeNode) genStats() TreeNodeStats {
	if tn.GetTierCode() == constants.TierMetaAgent.Code {
		return TreeNodeStats{node: tn}
	}

	// TODO: Add more stats fields here

	return TreeNodeStats{node: tn}
}

func (t *Tree) GetRoot() *TreeNode {
	return t.root
}

func (t *Tree) GetNode(userID meta.UID) *TreeNode {
	return t.nodeMap[userID]
}

func (t *Tree) GetNodeMap() TreeNodeMap {
	return t.nodeMap
}

func GenerateTreeAtNow() *Tree {
	return GenerateTreeAtDate("")
}

func GenerateTreeAtDate(date string) *Tree {
	root, nodeMap := generateTreeMap(date)

	return &Tree{
		root:    root,
		nodeMap: nodeMap,
	}
}

func generateTreeMap(date string) (*TreeNode, TreeNodeMap) {
	var allUsers []models.User
	comutils.PanicOnError(
		database.GetDbSlave().Find(&allUsers).Error,
	)

	var (
		nodeMap = make(TreeNodeMap, len(allUsers))
		userIDs = make([]meta.UID, 0, len(allUsers))
	)
	for _, user := range allUsers {
		nodeMap[user.ID] = genTreeNode(user)
		userIDs = append(userIDs, user.ID)
	}

	var userCoinBalanceMapMap map[meta.UID]tradingbalance.CoinBalanceMap
	if date == "" {
		userCoinBalanceMapMap = tradingbalance.GenUserCoinBalanceMapMapF(userIDs, true)
	} else {
		userCoinBalanceMapMap = tradingbalance.GenUserCoinStartBalanceMapMap(userIDs, true)
	}

	rootNode := &TreeNode{}
	for _, user := range allUsers {
		userNode := nodeMap[user.ID]
		userNode.CoinBalanceMap = userCoinBalanceMapMap[user.ID]

		parentNode, hasParent := nodeMap[userNode.User.ParentID]
		if !hasParent {
			if userNode.User.ParentID != RootUID {
				comlogging.GetLogger().
					WithFields(logrus.Fields{
						"uid":       user.ID,
						"parent_id": userNode.User.ParentID,
					}).
					Error("parent user not found")
				continue
			}

			parentNode = rootNode
		}

		userNode.Parent = parentNode
		parentNode.Children = append(parentNode.Children, userNode)
	}

	return rootNode, nodeMap
}

func genTreeNode(user models.User) *TreeNode {
	if user.IsDeleted.Bool || user.Status != "active" {
		dummyUser := models.User{
			ID:         user.ID,
			Code:       user.Code,
			FirstName:  "Torque User",
			ParentID:   user.ParentID,
			CreateDate: user.CreateDate,
		}
		user = dummyUser
	}

	treeUser := TreeUser{
		CreateTime: user.CreateDate.Unix(),
	}
	comutils.PanicOnError(
		copier.Copy(&treeUser, &user),
	)

	return &TreeNode{
		User: treeUser,
	}
}
