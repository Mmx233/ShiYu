package Service

type check struct {}
var Check check

func (*check)AccountExist(role string,username string)bool{
	var tempB bool
	DB.QueryRow("SELECT 1 FROM "+role+" WHERE username=?",username).Scan(&tempB)
	return tempB
}
