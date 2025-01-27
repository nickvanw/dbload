# sharded keyspace
CREATE TABLE `data` (
  `id` int NOT NULL,
  `host` varchar(100) DEFAULT NULL,
  `data` varchar(100) DEFAULT NULL,
  `now` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

# lookup keyspace
create table data_seq(id int, next_id bigint, cache bigint, primary key(id)) comment 'vitess_sequence';
insert into data_seq(id, next_id, cache) values(0, 1, 3);
