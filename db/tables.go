package db

import (
	"database/sql"
	"fmt"
	//"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type schemaCommand struct {
	description string
	create      string
}

var schemas = []schemaCommand{
	{
		description: "Project Schema",
		create: `
      create table projects (
           id integer primary key,
           name text,
           path text,
           mcid varchar(40)
      )`,
	},
	{
		description: "Project Changes Schema",
		create: `
      create table project_events (
           id integer primary key,
           path text,
           event varchar(40),
           event_time datetime,
           project_id integer,
           foreign key (project_id) references projects(id)     
      )`,
	},
	{
		description: "DataDir Schema",
		create: `
      create table datadirs (
           id integer primary key,
           project_id integer,
           mcid varchar(40),
           name text,
           path text,
           parent_mcid varchar(40),
           parent integer,
           foreign key (parent) references datadirs(id),
           foreign key (project_id) references projects(id)
      )`,
	},
	{
		description: "DataFile Schema",
		create: `
      create table datafiles (
           id integer primary key,
           mcid varchar(40),
           name text,
           path text,
           datadir_id integer,
           project_id integer,
           size integer,
           checksum varchar(64),
           last_upload datetime,
           mtime datetime,
           version integer,
           parent_mcid varchar(40),
           parent integer,
           foreign key (parent) references datafiles(id),
           foreign key (datadir_id) references datadirs(id),
           foreign key (project_id) references projects(id)
      )`,
	},
	{
		description: "Project To DataDir Join Table",
		create: `
      create table project2datadir (
           project_id integer,
           datadir_id integer,
           foreign key (project_id) references projects(id),
           foreign key (datadir_id) references datadirs(id)
      )`,
	},
	{
		description: "Project To DataFile Join Table",
		create: `
      create table project2datafile (
           project_id integer,
           datafile_id integer,
           foreign key (project_id) references projects(id),
           foreign key (datafile_id) references datafiles(id)
      )`,
	},
	{
		description: "DataDir To DataFile Join Table",
		create: `
      create table datadir2datafile (
           datadir_id integer,
           datafile_id integer,
           foreign key (datadir_id) references datadirs(id),
           foreign key (datafile_id) references datafiles(id)
      )`,
	},
	{
		description: "Trigger To Update Projects For DataDirs",
		create: `
  create trigger project_datadir after insert on datadirs
  begin
    insert into project2datadir(project_id, datadir_id) values(new.project_id, new.id);
  end`,
	},
	{
		description: "Trigger to Update Projects For DataFiles",
		create: `
      create trigger datafile_insert_updater after insert on datafiles
      begin
        insert into project2datafile(project_id, datafile_id) values(new.project_id, new.id);
        insert into datadir2datafile(datadir_id, datafile_id) values(new.datadir_id, new.id);
     end`,
	},
}

func Create(path string) error {
	dbArgs := fmt.Sprintf("file:%s?cached=shared&mode=rwc", path)
	db, err := sql.Open("sqlite3", dbArgs)
	if err != nil {
		return err
	}
	defer db.Close()

	for _, schema := range schemas {
		_, err = db.Exec(schema.create)
		if err != nil {
			return fmt.Errorf("Failed on create for %s: %s", schema.description, err)
		}
	}

	return nil
}
