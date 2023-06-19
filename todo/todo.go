package todo //ชื่อควรเหมือนกันหมด ณ ที่นี้คือ todo หมดเลย

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Todo struct {
	Title string `json:"text"`
	gorm.Model
}

// func (Todo) TableName() string {
// 	return "todos"
// }

func (Todo) TableName() string {
	return "to-do list" //ตั้งชื่อแล้วแต่เรา แบบนี้ก็ตั้งได้
}

type TodoHandler struct {
	db *gorm.DB //ดึง db มาจาก gorm แล้วตั้งชื่อว่า db (หรือจริงๆแล้วมันคือตั้งชื่อ db แล้วเขียน type เป็น gorm หรือเปล่า?)
}

func NewTodoHandler(db *gorm.DB) *TodoHandler {
	return &TodoHandler{db: db}
}

func (t *TodoHandler) NewTask(c *gin.Context) {
	var todo Todo // var var-name type
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}
	// ถ้า bind ข้อมูลแล้วไม่มีปัญหา มันจะไปทำงานในบรรทัดถัดไป

	if todo.Title == "sleep" {

		transactionID := c.Request.Header.Get("TransactionID")
		aud, _ := c.Get("aud")
		log.Println(transactionID, aud, "NOT allowed")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "NOT ALLOWED",
		})
	}

	// ทำการ create ของ แล้วยัดใส่ todo
	r := t.db.Create(&todo) // creaet จะคืน result มาตัวนึง ซึ่ง result จะมี error และถ้ามันไม่เท่ากับ nil ก็ handle กันไป
	if err := r.Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// อันนี้คือการคืนของไปให้หน้าบ้าน ใน gin.H คือจะคืนอะไรกลับไปบ้าง
	c.JSON(http.StatusCreated, gin.H{
		"ID": todo.Model.ID,
	})

}
