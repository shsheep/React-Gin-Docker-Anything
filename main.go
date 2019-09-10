package main

import (
   "fmt"
   "io"
   "io/ioutil"
   "os"
   "context"
   "path/filepath"
   "encoding/json"
   "time"

   /* Server related packages */
   "github.com/gin-gonic/gin"
   "github.com/gin-gonic/contrib/static"
   "net/http"

   /* Docker client related pakcages */
   "github.com/docker/docker/api/types"
   "github.com/docker/docker/api/types/container"
   "github.com/docker/docker/client"
   "github.com/docker/docker/api/types/mount"
   //"github.com/docker/docker/api/types/volume"

   /* User package for managing SR containers */
   "srcontainer"
)

/* DEBUGGING FLAG - A flag which determines whether to print docker-related, server-related logs or not */
var IGNITION_LOG_FLAG bool = true

func main() {
   /* * * * * * * * * * *  Initial variables related to containers * * * * * * * * * * * * * */
   Media_image_name := "localhost:5000/shsheep/media_with_script"
   TAB_image_name := "localhost:5000/shsheep/tab_with_script"
   media_input_repository := "/home/shsheep/Workspace/React-Gin-Docker-Anything/media_input_repository"
   media_output_repository := "/home/shsheep/Workspace/React-Gin-Docker-Anything/media_output_repository"
   tab_input_repository := "/home/shsheep/Workspace/React-Gin-Docker-Anything/tab_input_repository"
   tab_output_repository := "/home/shsheep/Workspace/React-Gin-Docker-Anything/tab_output_repository"
   /* * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * * *  */

   ctx := context.Background()

   /* Initiailize the client object with specific version */
   cli, err := client.NewClientWithOpts(client.WithVersion("1.39"))
   if err != nil {
      panic(err)
   }

   /* Check currently running containers */
   if IGNITION_LOG_FLAG {
      srcontainer.CheckRunningContainers(cli, ctx)
   }
   container_list := srcontainer.AssignContainerNumber()
   tab_reader := srcontainer.PullImage(cli, ctx, TAB_image_name)
   media_reader := srcontainer.PullImage(cli, ctx, Media_image_name)

   /* Check whether the image is normally pulled */
   if IGNITION_LOG_FLAG {
      io.Copy(os.Stdout, media_reader)
      io.Copy(os.Stdout, tab_reader)
   }

   /* Create each containers */
   Media_resp, err := cli.ContainerCreate(ctx,
      &container.Config{
         Image: Media_image_name,
         Tty: true,
         },
      &container.HostConfig{
         Mounts: []mount.Mount {
            {
               Type: mount.TypeBind,
               Source: media_input_repository,
               Target: "/media/input",
            },
            {
               Type: mount.TypeBind,
               Source: media_output_repository,
               Target: "/media/output",
            },
         },
      }, nil, container_list[0])
   if err != nil {
      panic(err)
   }

   Tab_resp, err := cli.ContainerCreate(ctx,
      &container.Config{
         Image: TAB_image_name,
         Tty: true,
         },
      &container.HostConfig{
         Runtime: "nvidia",
         Mounts: []mount.Mount {
            {
               Type: mount.TypeBind,
               Source: tab_output_repository,
               Target: "/tab/output",
            },
            {
               Type: mount.TypeBind,
               Source: tab_input_repository,
               Target: "/tab/input",
            },
         },
      }, nil, container_list[1])
   if err != nil {
      panic(err)
   }

   /* Start(Run) the container */
   if err := cli.ContainerStart(ctx, Media_resp.ID, types.ContainerStartOptions{}); err != nil {
      panic(err)
   }
   if err := cli.ContainerStart(ctx, Tab_resp.ID, types.ContainerStartOptions{}); err != nil {
      panic(err)
   }

   /* Check the container logs */
   if IGNITION_LOG_FLAG {
      srcontainer.CheckLogs(cli, ctx, Media_resp)
      srcontainer.CheckLogs(cli, ctx, Tab_resp)
   }

   /* . . . . . . GIN-GONIC SERVER RUNS . . . . . . */
   router := gin.Default()
   router.Use(static.Serve("/", static.LocalFile("./client/build", true)))
   api := router.Group("/api")
   {
      byteTestResult := make(map[string]interface{})
      bytePostResult := make(map[string]interface{})

   /* * * * * * * * * * * * * * * * * * MEDIA * * * * * * * * * * * * * * * * * * * */
      /* CASE : when Media input file upload */
      api.POST("/media-upload", func(c *gin.Context) {
         now := time.Now()
         file, err := c.FormFile("file")
         if err != nil {
            panic(err)
         }

         /* Name the file with timestamp */
         filename := filepath.Base(file.Filename) + now.Format("2006-01-02T15:04:05")
         uploadPath := media_input_repository + "/" + filename

         if err := c.SaveUploadedFile(file, uploadPath); err != nil {
            panic(err)
         }
         bytePostResult["filename"] = filename
         c.JSON(http.StatusOK, gin.H{
            "filename": byteTestResult["filename"],
         })
      })

      /* RESULT PARSING */
      api.GET("/media-get-result", func(c *gin.Context) {
         result, err := ioutil.ReadFile(media_output_repository + "/result.out")
         if err != nil {
            panic(err)
         }
         byteTestResult["result"] = string(result)
         c.JSON(http.StatusOK, gin.H{
            "result": byteTestResult["result"],
         })
      })

   /* * * * * * * * * * * * * * * * * * TAB * * * * * * * * * * * * * * * * * * * */
      /* CASE : for input file upload */
      api.POST("/tab-upload", func(c *gin.Context) {
         now := time.Now()
         file, err := c.FormFile("file")
         if err != nil {
            panic(err)
         }

         /* Name the file with timestamp */
         filename := filepath.Base(file.Filename) + now.Format("2006-01-02T15:04:05")
         uploadPath := tab_input_repository + "/" + filename
         if err := c.SaveUploadedFile(file, uploadPath); err != nil {
            panic(err)
         }
      })

      /* CASE : for manually typed text */
      api.POST("/tab-write-get-result", func(c *gin.Context) {
         now := time.Now()
         var byteTAB map[string]interface{}
         tmp, _ := c.GetRawData()

         /* Name the file with timestamp */
         json.Unmarshal(tmp, &byteTAB)
         newFileName := "Manual_input" + now.Format("2006-01-02T15:04:05")
         file, _ := os.Create(filepath.Join(tab_input_repository, filepath.Base(newFileName)))
         defer file.Close()

         fmt.Fprint(file, byteTAB["content"])
      })

      /* RESULT PARSING */
      api.GET("/tab-get-result", func(c *gin.Context) {
         result, err := ioutil.ReadFile(tab_output_repository + "/json_output.txt")
         if err != nil {
            panic(err)
         }

         byteTestResult["result"] = string(result)
         c.JSON(http.StatusOK, gin.H{
            "result": byteTestResult["result"],
         })
      })
   }
   router.Run(":8000")
}
