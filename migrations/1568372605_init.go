package migrations

import migrate "github.com/rubenv/sql-migrate"

func init() {
	instance.add(&migrate.Migration{
		Id: "1568372605",
		Up: []string{
			`
            CREATE TABLE category (
			  id bigint(20) NOT NULL AUTO_INCREMENT,
			  created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
			  updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			  deleted_at timestamp NULL DEFAULT NULL,
			  category_name varchar(255) NOT NULL,
			  PRIMARY KEY (id),
			  UNIQUE KEY id (id),
			  UNIQUE KEY category_name (category_name)
			);
			`,
			`
			CREATE TABLE product (
			  id bigint(20) NOT NULL AUTO_INCREMENT,
			  created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
			  updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
			  deleted_at timestamp NULL DEFAULT NULL,
			  data_asin varchar(255) NOT NULL,
			  product_name varchar(255) NULL,
			  manufacturer varchar(255) NULL,
			  category_id bigint(20) NULL,
			  price varchar(255) NULL,
			  ratings varchar(255) NULL,
			  description TEXT NULL,
			  PRIMARY KEY (id),
			  UNIQUE KEY id (id),
			  UNIQUE KEY data_asin (data_asin),
			  FOREIGN KEY (category_id) REFERENCES category (id) ON DELETE CASCADE
			);
           `,
		},
		//language=SQL
		Down: []string{
			`DROP TABLE product;`,
			`DROP TABLE category;`,
		},
	})
}
