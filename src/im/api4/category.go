package api4

import (
	"fmt"
	"im/model"
	"net/http"
)

func (api *API) InitCategory() {
	api.BaseRoutes.Categories.Handle("", api.ApiHandler(getCategories)).Methods("GET")
	api.BaseRoutes.Categories.Handle("", api.ApiHandler(createCategory)).Methods("POST")

	api.BaseRoutes.Category.Handle("", api.ApiHandler(getCategory)).Methods("GET")
	api.BaseRoutes.Category.Handle("", api.ApiHandler(updateCategory)).Methods("PUT")
	api.BaseRoutes.Category.Handle("", api.ApiHandler(deleteCategory)).Methods("DELETE")
}

func getCategory(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireCategoryId()
	if c.Err != nil {
		return
	}
	category, err := c.App.GetCategory(c.Params.CategoryId)
	if err != nil {
		c.Err = err
		return
	}

	w.Write([]byte(category.ToJson()))
}

func getCategories(c *Context, w http.ResponseWriter, r *http.Request) {
	categories, err := c.App.GetCategoriesPage(0, c.Params.PerPage)
	if err != nil {
		c.Err = err
		return
	}
	w.Write([]byte(model.CategoriesToJson(categories)))
}

func getCategoriesByClient(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireClientId()
	if c.Err != nil {
		return
	}
	categories, err := c.App.GetCategoriesByClientIdPage(c.Params.ClientId, c.Params.Page, c.Params.PerPage)
	if err != nil {
		c.Err = err
		return
	}
	w.Write([]byte(model.CategoriesToJson(categories)))
}

func createCategory(c *Context, w http.ResponseWriter, r *http.Request) {
	category := model.CategoryFromJson(r.Body)
	if category == nil {
		c.Err = model.NewAppError("createCategory", "api.category", nil, "nil object", http.StatusForbidden)
		return
	}
	result, err := c.App.CreateCategory(category)
	if err != nil {
		c.Err = err
		return
	}
	w.Write([]byte(result.ToJson()))
}

func updateCategory(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireCategoryId()
	if c.Err != nil {
		return
	}

	category := model.CategoryFromJson(r.Body)
	fmt.Print(category)
	if category == nil {
		c.SetInvalidParam("category")
		return
	}

	// The post being updated in the payload must be the same one as indicated in the URL.
	if category.Id != c.Params.CategoryId {
		c.SetInvalidParam("id")
		return
	}

	// Updating the file_ids of a post is not a supported operation and will be ignored

	/*	if !c.App.SessionHasPermissionToChannelByPost(c.App.Session, c.Params.PostId, model.PERMISSION_EDIT_POST) {
		c.SetPermissionError(model.PERMISSION_EDIT_POST)
		return
	}*/

	/*originalCategory, err := c.App.GetSingleCategory(c.Params.CategoryId)
	if err != nil {
		c.SetPermissionError(model.PERMISSION_EDIT_POST)
		return
	}*/

	category.Id = c.Params.CategoryId

	rcategory, err := c.App.UpdateCategory(category, false)
	if err != nil {
		c.Err = err
		return
	}

	w.Write([]byte(rcategory.ToJson()))
}

func _deleteCategory(c *Context, w http.ResponseWriter, r *http.Request) {
	c.RequireCategoryId()
	if c.Err != nil {
		return
	}

	category, err := c.App.GetCategory(c.Params.CategoryId)
	if err != nil {
		c.Err = err
		return
	}
	// The post being updated in the payload must be the same one as indicated in the URL.
	if category.Id != c.Params.CategoryId {
		c.SetInvalidParam("id")
		return
	}

	if _, err := c.App.DeleteCategory(category); err != nil {
		c.Err = err
		return
	}

	ReturnStatusOK(w)
}

func deleteCategory(c *Context, w http.ResponseWriter, r *http.Request) {
	category, err := c.App.GetCategory(c.Params.CategoryId)
	if err != nil {
		c.Err = err
		return
	}
	c.App.DeleteCategory(category)
	/*
		if err != nil {
			c.Err = err
			ReturnStatusOK(w)
		}
	*/
	ReturnStatusOK(w)

	//w.Write([]byte(model.MapToJson(map[string]string{"status": strconv.Itoa(result["status"])})))
}
