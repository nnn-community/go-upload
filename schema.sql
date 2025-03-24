-- Table Definition
CREATE TABLE "public"."uploaded_files" (
    "id" varchar(36) NOT NULL,
    "user_id" varchar(36),
    "filename" text NOT NULL,
    "size" int4 NOT NULL,
    "uploaded_url" text NOT NULL,
    "created_at" timestamp(3) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY ("id")
);

-- Indices
CREATE INDEX uploaded_files_id_idx ON public.uploaded_files USING hash (id);
CREATE INDEX uploaded_files_user_id_idx ON public.uploaded_files USING hash (user_id);
CREATE INDEX uploaded_files_filename_idx ON public.uploaded_files USING hash (filename);
CREATE INDEX uploaded_files_size_idx ON public.uploaded_files USING hash (size);
CREATE INDEX uploaded_files_created_at_idx ON public.uploaded_files USING hash (created_at);
