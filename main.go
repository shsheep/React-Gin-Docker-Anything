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
   MSD_image_name := "localhost:5000/shsheep/msd_with_script"
   SLU_image_name := "localhost:5000/shsheep/slu_with_script"
   msd_input_repository := "/home/shsheep/Workspace/React-Gin-Docker-Anything/msd_input_repository"
   msd_output_repository := "/home/shsheep/Workspace/React-Gin-Docker-Anything/msd_output_repository"
   slu_input_repository := "/home/shsheep/Workspace/React-Gin-Docker-Anything/slu_input_repository"
   slu_output_repository := "/home/shsheep/Workspace/React-Gin-Docker-Anything/slu_output_repository"
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
   slu_reader := srcontainer.PullImage(cli, ctx, SLU_image_name)
   msd_reader := srcontainer.PullImage(cli, ctx, MSD_image_name)

   /* Check whether the image is normally pulled */
   if IGNITION_LOG_FLAG {
      io.Copy(os.Stdout, msd_reader)
      io.Copy(os.Stdout, slu_reader)
   }

   /* Create each containers */
   MSD_resp, err := cli.ContainerCreate(ctx,
      &container.Config{
         Image: MSD_image_name,
         Tty: true,
         },
      &container.HostConfig{
         Mounts: []mount.Mount {
            {
               Type: mount.TypeBind,
               Source: msd_input_repository,
               Target: "/msd/input",
            },
            {
               Type: mount.TypeBind,
               Source: msd_output_repository,
               Target: "/msd/output",
            },
         },
      }, nil, container_list[0])
   if err != nil {
      panic(err)
   }

   SLU_resp, err := cli.ContainerCreate(ctx,
      &container.Config{
         Image: SLU_image_name,
         Tty: true,
         },
      &container.HostConfig{
         Runtime: "nvidia",
         Mounts: []mount.Mount {
            {
               Type: mount.TypeBind,
               Source: slu_output_repository,
               Target: "/slu/output",
            },
            {
               Type: mount.TypeBind,
               Source: slu_input_repository,
               Target: "/slu/input",
            },
         },
      }, nil, container_list[1])
   if err != nil {
      panic(err)
   }

   /* Start(Run) the container */
   if err := cli.ContainerStart(ctx, MSD_resp.ID, types.ContainerStartOptions{}); err != nil {
      panic(err)
   }
   if err := cli.ContainerStart(ctx, SLU_resp.ID, types.ContainerStartOptions{}); err != nil {
      panic(err)
   }

   /* Check the container logs */
   if IGNITION_LOG_FLAG {
      srcontainer.CheckLogs(cli, ctx, MSD_resp)
      srcontainer.CheckLogs(cli, ctx, SLU_resp)
   }

   /* . . . . . . GIN-GONIC SERVER RUNS . . . . . . */
   router := gin.Default()
   router.Use(static.Serve("/", static.LocalFile("./client/build", true)))
   api := router.Group("/api")
   {
      byteTestResult := make(map[string]interface{})
      bytePostResult := make(map[string]interface{})

   /* * * * * * * * * * * * * * * * * * MSD * * * * * * * * * * * * * * * * * * * */
      /* CASE : when MSD input file upload */
      api.POST("/msd-upload", func(c *gin.Context) {
         now := time.Now()
         file, err := c.FormFile("file")
         if err != nil {
            panic(err)
         }

         /* Name the file with timestamp */
         filename := filepath.Base(file.Filename) + now.Format("2006-01-02T15:04:05")
         uploadPath := msd_input_repository + "/" + filename

         if err := c.SaveUploadedFile(file, uploadPath); err != nil {
            panic(err)
         }
         bytePostResult["filename"] = filename
         c.JSON(http.StatusOK, gin.H{
            "filename": byteTestResult["filename"],
         })
      })

      /* RESULT PARSING */
      api.GET("/msd-get-result", func(c *gin.Context) {
         result, err := ioutil.ReadFile(msd_output_repository + "/result.out")
         if err != nil {
            panic(err)
         }
         byteTestResult["result"] = string(result)
         c.JSON(http.StatusOK, gin.H{
            "result": byteTestResult["result"],
         })
      })

   /* * * * * * * * * * * * * * * * * * SLU * * * * * * * * * * * * * * * * * * * */
      /* CASE : for input file upload */
      api.POST("/slu-upload", func(c *gin.Context) {
         now := time.Now()
         file, err := c.FormFile("file")
         if err != nil {
            panic(err)
         }

         /* Name the file with timestamp */
         filename := filepath.Base(file.Filename) + now.Format("2006-01-02T15:04:05")
         uploadPath := slu_input_repository + "/" + filename
         if err := c.SaveUploadedFile(file, uploadPath); err != nil {
            panic(err)
         }
      })

      /* CASE : for manually typed text */
      api.POST("/slu-write-get-result", func(c *gin.Context) {
         now := time.Now()
         var byteSLU map[string]interface{}
         tmp, _ := c.GetRawData()

         /* Name the file with timestamp */
         json.Unmarshal(tmp, &byteSLU)
         newFileName := "Manual_input" + now.Format("2006-01-02T15:04:05")
         file, _ := os.Create(filepath.Join(slu_input_repository, filepath.Base(newFileName)))
         defer file.Close()

         fmt.Fprint(file, byteSLU["content"])
      })

      /* RESULT PARSING */
      api.GET("/slu-get-result", func(c *gin.Context) {
         result, err := ioutil.ReadFile(slu_output_repository + "/json_output.txt")
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
