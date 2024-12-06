package controllers

type Controller struct {
	Auth   *AuthController
	Tender *TenderController
	Offer  *OfferController
	Notif  *NotifController
}
