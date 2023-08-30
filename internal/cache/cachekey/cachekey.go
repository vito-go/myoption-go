package cachekey

const (
	UserFolloweeList = "myoption:list:followees:%s"   // userId
	UserBlackListing = "myoption:list:blacklisted:%s" // userId
)
const (
	KeyGeneralDistributeDoOnce = `myoption:KeyGeneralDistributeDoOnce:%s` // flag
	KeyGeneralDistributeLock   = `myoption:KeyGeneralDistributeLock:%s`   // flag
)

const (
	UserInfo       = "myoption:userInfo:%s"              // useId
	UserKey        = "myoption:UserKey:%s"               // useId
	UserLogin      = "myoption:user-login:%s"            // useId
	UserLoginHash  = "myoption:user-login-hash:%s"       // useId //多端登陆准备
	UserLoginToken = "myoption:user-loginInfo:%s"        // 未未來多端登錄支持做準備
	LastOnline     = "myoption:user-lastOnline:%s"       // useId
	UserBlackList  = "myoption:user-black-list-map:%s"   // userId
	UserFollower   = "myoption:user-follower-map:%s"     // userId
	UserFollowee   = "myoption:user-followee-map:%s"     // userId
	UserLikeMoment = "myoption:hash:user_like_moment:%s" // userId
)
