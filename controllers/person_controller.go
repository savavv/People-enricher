package controllers

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"people-enricher/database"
	"people-enricher/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetPeople возвращает список людей с фильтрами и пагинацией.
// @Summary Получить список людей
// @Description Поддерживает фильтры по имени, фамилии и пагинацию
// @Tags people
// @Accept json
// @Produce json
// @Param name query string false "Имя"
// @Param surname query string false "Фамилия"
// @Param page query int false "Номер страницы" default(1)
// @Param limit query int false "Размер страницы" default(10)
// @Success 200 {array} models.Person
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /people [get]
func GetPeople(c *gin.Context) {
	var people []models.Person

	query := database.DB.Model(&models.Person{})

	if name := c.Query("name"); name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}
	if surname := c.Query("surname"); surname != "" {
		query = query.Where("surname ILIKE ?", "%"+surname+"%")
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	offset := (page - 1) * limit

	query = query.Offset(offset).Limit(limit)

	if err := query.Find(&people).Error; err != nil {
		log.Printf("[ERROR] failed to query people: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch people"})
		return
	}

	log.Printf("[INFO] fetched %d people (page %d, limit %d)", len(people), page, limit)
	c.JSON(http.StatusOK, people)
}

// AddPerson добавляет нового человека в базу.
// @Summary Добавить нового человека
// @Description Принимает имя, фамилию и отчество (необязательно), обогащает данными и сохраняет
// @Tags people
// @Accept json
// @Produce json
// @Param person body models.Person true "Данные человека"
// @Success 201 {object} models.Person
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /people [post]
func AddPerson(c *gin.Context) {
	var person models.Person

	if err := c.ShouldBindJSON(&person); err != nil {
		log.Printf("[ERROR] invalid request body: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	log.Printf("[INFO] Received new person: %+v", person)

	// Обогащаем
	person.Age = getAge(person.Name)
	person.Gender = getGender(person.Name)
	person.Nationality = getNationality(person.Name)

	log.Printf("[DEBUG] Enriched person: %+v", person)

	if err := database.DB.Create(&person).Error; err != nil {
		log.Printf("[ERROR] failed to save person: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save person"})
		return
	}

	c.JSON(http.StatusCreated, person)
}

// UpdatePerson обновляет данные человека по ID.
// @Summary Обновить человека по ID
// @Description Обновляет данные существующего человека и обогащает их
// @Tags people
// @Accept json
// @Produce json
// @Param id path int true "ID человека"
// @Param person body models.Person true "Обновленные данные человека"
// @Success 200 {object} models.Person
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /people/{id} [put]
func UpdatePerson(c *gin.Context) {
	id := c.Param("id")
	var person models.Person

	if err := database.DB.First(&person, id).Error; err != nil {
		log.Printf("[ERROR] person not found: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Person not found"})
		return
	}

	var input models.Person
	if err := c.ShouldBindJSON(&input); err != nil {
		log.Printf("[ERROR] invalid update data: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	person.Name = input.Name
	person.Surname = input.Surname
	person.Patronymic = input.Patronymic

	// Обогащаем заново
	person.Age = getAge(person.Name)
	person.Gender = getGender(person.Name)
	person.Nationality = getNationality(person.Name)

	log.Printf("[DEBUG] Updated person: %+v", person)

	if err := database.DB.Save(&person).Error; err != nil {
		log.Printf("[ERROR] failed to update person: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update person"})
		return
	}

	c.JSON(http.StatusOK, person)
}

// DeletePerson удаляет человека по ID.
// @Summary Удалить человека по ID
// @Description Удаляет человека из базы
// @Tags people
// @Accept json
// @Produce json
// @Param id path int true "ID человека"
// @Success 200 {object} map[string]string "Person deleted"
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /people/{id} [delete]
func DeletePerson(c *gin.Context) {
	id := c.Param("id")

	if err := database.DB.Delete(&models.Person{}, id).Error; err != nil {
		log.Printf("[ERROR] failed to delete person: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete person"})
		return
	}

	log.Printf("[INFO] Deleted person ID: %s", id)
	c.JSON(http.StatusOK, gin.H{"message": "Person deleted"})
}

func getAge(name string) int {
	url := "https://api.agify.io/?name=" + name
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("[ERROR] failed to get age: %v", err)
		return 0
	}
	defer resp.Body.Close()

	var res struct {
		Age int `json:"age"`
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ERROR] failed to read age response: %v", err)
		return 0
	}

	if err := json.Unmarshal(body, &res); err != nil {
		log.Printf("[ERROR] failed to unmarshal age response: %v", err)
		return 0
	}

	log.Printf("[DEBUG] Age from agify.io for name %s: %d", name, res.Age)
	return res.Age
}

func getGender(name string) string {
	url := "https://api.genderize.io/?name=" + name
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("[ERROR] failed to get gender: %v", err)
		return ""
	}
	defer resp.Body.Close()

	var res struct {
		Gender string `json:"gender"`
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ERROR] failed to read gender response: %v", err)
		return ""
	}

	if err := json.Unmarshal(body, &res); err != nil {
		log.Printf("[ERROR] failed to unmarshal gender response: %v", err)
		return ""
	}

	log.Printf("[DEBUG] Gender from genderize.io for name %s: %s", name, res.Gender)
	return res.Gender
}

func getNationality(name string) string {
	url := "https://api.nationalize.io/?name=" + name
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("[ERROR] failed to get nationality: %v", err)
		return ""
	}
	defer resp.Body.Close()

	var res struct {
		Country []struct {
			CountryID   string  `json:"country_id"`
			Probability float64 `json:"probability"`
		} `json:"country"`
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[ERROR] failed to read nationality response: %v", err)
		return ""
	}

	if err := json.Unmarshal(body, &res); err != nil {
		log.Printf("[ERROR] failed to unmarshal nationality response: %v", err)
		return ""
	}

	if len(res.Country) > 0 {
		log.Printf("[DEBUG] Nationality from nationalize.io for name %s: %s", name, res.Country[0].CountryID)
		return res.Country[0].CountryID
	}

	return ""
}
