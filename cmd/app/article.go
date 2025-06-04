package app

import (
	"fmt"
	"math"

	"errors"

	"github.com/yann0917/dedao-dl/services"
	"github.com/yann0917/dedao-dl/utils"
)

// ArticleList 已购课程文章列表
func ArticleList(id int, chapterID string) (list *services.ArticleList, err error) {
	info, err := CourseDetail(CateCourse, id)
	if err != nil {
		return
	}
	enid := info.Enid
	count := float64(info.PublishNum)
	page := int(math.Ceil(float64(count) / 20.0))
	maxID := 0
	var lists []services.ArticleIntro
	for i := 0; i < page; i++ {
		list, err = getService().ArticleList(enid, chapterID, maxID)
		if err != nil {
			return
		}
		if len(list.List) > 0 {
			maxID = list.List[len(list.List)-1].ID
		}
		lists = append(lists, list.List...)
	}
	list.List = lists
	return

}

// ArticleInfo article info
// get article token, audio token, media security token etc
func ArticleInfo(id, aid int) (info *services.ArticleInfo, aEnid string, err error) {
	list, err := ArticleList(id, "")
	if err != nil {
		return
	}

	var aids []int

	// get article enid
	for _, p := range list.List {
		aids = append(aids, p.ID)
		if p.ClassID == id && p.ID == aid {
			aEnid = p.Enid
		}
	}
	if !utils.Contains(aids, aid) {
		err = errors.New("找不到该文章 ID，请检查输入是否正确")
		return
	}

	info, err = getService().ArticleInfo(aEnid, 1)
	if err != nil {
		return
	}
	return
}

// ArticleDetail article detail
func ArticleDetail(id, aid int) (detail *services.ArticleDetail, enId string, err error) {
	info, aEnid, err := ArticleInfo(id, aid)
	if err != nil {
		return
	}
	enId = aEnid
	token := info.DdArticleToken
	appid := "1632426125495894021"
	detail, err = getService().ArticleDetail(token, aEnid, appid)
	if err != nil {
		fmt.Printf("err:%#v\n", err)
		return
	}
	return

}

// ArticleCommentList article comment list
func ArticleCommentList(enId, sort string, page, limit int) (list *services.ArticleCommentList, err error) {
	list, err = getService().ArticleCommentList(enId, sort, page, limit)
	if err != nil {
		return
	}

	return

}

// OdobArticleInfo article info
// get article token, audio token, media security token etc.
func OdobArticleInfo(aEnid string) (info *services.ArticleInfo, err error) {
	info, err = getService().ArticleInfo(aEnid, 2)
	if err != nil {
		return
	}
	return
}

// OdobArticleDetail odob article detail
func OdobArticleDetail(aEnid string) (detail *services.ArticleDetail, err error) {
	info, err := OdobArticleInfo(aEnid)
	if err != nil {
		return
	}

	token := info.DdArticleToken
	appid := "1632426125495894021"
	detail, err = getService().ArticleDetail(token, aEnid, appid)
	if err != nil {
		fmt.Printf("err:%#v\n", err)
		return
	}
	return

}
