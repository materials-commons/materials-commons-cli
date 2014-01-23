package db

import (
	"github.com/jmoiron/sqlx"
	"github.com/materials-commons/materials/db/model"
	"github.com/materials-commons/materials/db/schema"
)

var (
	ProjectsModel      *model.Model
	ProjectEventsModel *model.Model
	DataDirsModel      *model.Model
	DataFilesModel     *model.Model
	Projects           *model.Query
	ProjectEvents      *model.Query
	DataDirs           *model.Query
	DataFiles          *model.Query
)

func Use(db *sqlx.DB) {
	Projects = ProjectsModel.Q(db)
	ProjectEvents = ProjectEventsModel.Q(db)
	DataDirs = DataDirsModel.Q(db)
	DataFiles = DataFilesModel.Q(db)
}

func init() {
	pQueries := model.ModelQueries{
		Insert: "insert into projects (name, path, mcid) values (:name, :path, :mcid)",
	}
	ProjectsModel = model.New(schema.Project{}, "projects", pQueries)

	peQueries := model.ModelQueries{
		Insert: `insert into project_events (path, event, event_time, project_id)
                 values (:path, :event, :event_time, :project_id)`,
	}
	ProjectEventsModel = model.New(schema.ProjectEvent{}, "project_events", peQueries)

	ddirQueries := model.ModelQueries{
		Insert: `insert into datadirs (mcid, project_id, name, path, parent_mcid, parent)
                 values (:mcid, :project_id, :name, :path, :parent_mcid, :parent)`,
	}
	DataDirsModel = model.New(schema.DataDir{}, "datadirs", ddirQueries)

	dfQueries := model.ModelQueries{
		Insert: `insert into datafiles
                  (mcid, name, path, datadir_id, project_id, size, checksum, last_upload, mtime, version, parent_mcid, parent)
                 values
                  (:mcid, :name, :path, :datadir_id, :project_id, :size, :checksum, :last_upload, :mtime, :version, :parent_mcid, :parent)`,
	}
	DataFilesModel = model.New(schema.DataFile{}, "datafiles", dfQueries)
}
