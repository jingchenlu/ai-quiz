-- ----------------------------
-- Table structure for paper_questions
-- ----------------------------
CREATE TABLE IF NOT EXISTS "paper_questions" (
    "id" integer PRIMARY KEY AUTOINCREMENT,
    "paper_id" integer NOT NULL,
    "question_id" integer NOT NULL,
    "question_order" integer NOT NULL,
    "score" integer DEFAULT 5,
    "created_at" datetime,
    "updated_at" datetime,
    "deleted_at" datetime,
    CONSTRAINT "fk_paper_questions_question" FOREIGN KEY ("question_id") REFERENCES "questions" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION,
    CONSTRAINT "fk_papers_questions" FOREIGN KEY ("paper_id") REFERENCES "papers" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- ----------------------------
-- Indexes for paper_questions
-- ----------------------------
-- 试卷题目顺序唯一索引（仅未删除记录生效）
CREATE UNIQUE INDEX IF NOT EXISTS "idx_paper_order_active"
    ON "paper_questions" ("paper_id" ASC, "question_order" ASC)
    WHERE deleted_at IS NULL;

-- 试卷与题目关联唯一索引（仅未删除记录生效）
CREATE UNIQUE INDEX IF NOT EXISTS "idx_paper_question_unique_active"
    ON "paper_questions" ("paper_id" ASC, "question_id" ASC)
    WHERE deleted_at IS NULL;

-- ----------------------------
-- Table structure for papers
-- ----------------------------
CREATE TABLE IF NOT EXISTS "papers" (
    "id" integer PRIMARY KEY AUTOINCREMENT,
    "title" text NOT NULL,
    "description" text,
    "total_score" integer DEFAULT 100,
    "creator_id" integer NOT NULL,
    "created_at" datetime,
    "updated_at" datetime,
    "deleted_at" datetime,
    CONSTRAINT "fk_papers_creator" FOREIGN KEY ("creator_id") REFERENCES "users" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- ----------------------------
-- Table structure for questions
-- ----------------------------
CREATE TABLE IF NOT EXISTS "questions" (
    "id" integer PRIMARY KEY AUTOINCREMENT,
    "title" text NOT NULL,
    "question_type" text NOT NULL,
    "options" text NOT NULL,
    "answer" text NOT NULL,
    "explanation" text,
    "keywords" text,
    "language" text NOT NULL,
    "ai_model" text NOT NULL,
    "user_id" integer NOT NULL,
    "created_at" datetime,
    "updated_at" datetime,
    "deleted_at" datetime,
    CONSTRAINT "fk_questions_user" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE NO ACTION ON UPDATE NO ACTION
);

-- ----------------------------
-- Table structure for users
-- ----------------------------
CREATE TABLE IF NOT EXISTS "users" (
    "id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,
    "username" text NOT NULL,
    "password_hash" text NOT NULL,
    "role" text DEFAULT "user",
    "created_at" datetime,
    "updated_at" datetime,
    "deleted_at" datetime,
    CONSTRAINT "uni_users_username" UNIQUE ("username" ASC)
);

-- ----------------------------
-- Indexes
-- ----------------------------
CREATE INDEX IF NOT EXISTS "idx_users_deleted_created"
    ON "users" ("deleted_at" ASC, "created_at" DESC);

CREATE INDEX IF NOT EXISTS "idx_questions_deleted_created"
    ON "questions" ("deleted_at" ASC, "created_at" DESC);

CREATE INDEX IF NOT EXISTS "idx_papers_deleted_created"
    ON "papers" ("deleted_at" ASC, "created_at" DESC);

-- ----------------------------
-- Initialize users table with admin account
-- ----------------------------
INSERT OR IGNORE INTO "users" VALUES (1, 'admin', '$2a$10$HnrCyQFMFY3aTB7kbsz//OsHpa.YLK172BVqbLYBJIiOJ0YNIachu', 'admin', '2025-07-21 09:19:32.138162+08:00', '2025-07-21 09:19:32.138162+08:00', NULL);