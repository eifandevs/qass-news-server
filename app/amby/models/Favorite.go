package models

import (
    "github.com/eifandevs/amby/repo"
    "github.com/thoas/go-funk"
)

type FavoriteItem struct {
    ID int `json:"id"`
    Title string `json:"title"`
    Url string `json:"url"`
}

type Favorite struct {
    ID int `gorm:"type:int unsigned;not null;primary_key;auto_increment:false"`
    Token string `gorm:"type:varchar(255);not null;primary_key"`
    Title string
    Url string
}

type GetFavoriteResponse struct {
    BaseResponse
    Items  []FavoriteItem `json:"data"`
}

type PostFavoriteRequest struct {
    Items  []FavoriteItem `json:"data"`
}

type DeleteFavoriteRequest struct {
    Items  []FavoriteItem `json:"data"`
}

func GetFavorite(userToken string) GetFavoriteResponse {
    db := repo.Connect("development")
    defer db.Close()

    favorites := []Favorite{}
    if err := db.Where("token = ?", userToken).Find(&favorites).Error; err != nil {
        return GetFavoriteResponse{BaseResponse: BaseResponse{Result: "NG", ErrorCode: ""}, Items: nil}
    }

    items := funk.Map(favorites, func(favorite Favorite) FavoriteItem {
        return FavoriteItem{ID: favorite.ID, Title: favorite.Title, Url: favorite.Url}
    })
    
    if castedItems, ok := items.([]FavoriteItem); ok {
        return GetFavoriteResponse{BaseResponse: BaseResponse{Result: "OK", ErrorCode: ""}, Items: castedItems}
    } else {
        panic("cannot cast favorite item.")
    }
}

func PostFavorite(userToken string, request PostFavoriteRequest) BaseResponse {
    db := repo.Connect("development")
    defer db.Close()

    for _, item := range request.Items {
        if err := db.Create(&Favorite{ID: item.ID, Token: userToken, Title: item.Title, Url: item.Url}).Error; err != nil {
            return BaseResponse{Result: "NG", ErrorCode: ""}
        }
    }

	return BaseResponse{Result: "OK", ErrorCode: ""}
}

func DeleteFavorite(userToken string, request DeleteFavoriteRequest) BaseResponse {
    db := repo.Connect("development")
    defer db.Close()

    for _, item := range request.Items {
        deletingRecord := Favorite{}
        deletingRecord.ID = item.ID
        deletingRecord.Token = userToken
        db.First(&deletingRecord)
        if err := db.Unscoped().Delete(&deletingRecord).Error; err != nil {
            return BaseResponse{Result: "NG", ErrorCode: ""}
        }
    }

	return BaseResponse{Result: "OK", ErrorCode: ""}
}