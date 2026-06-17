-- queries for video table 

INSERT INTO videos (id, filename, original_path, filesize, status) VALUES (?, ?, ?, ?, 'pending') RETURNING *;

SELECT * FROM videos WHERE id = ? LIMIT 1;

SELECT * FROM videos ORDER BY created_at DESC;

UPDATE videos SET duration = ?, width = ?, height = ?, fps = ?, codec = ?, status = 'ready', updated_at = strftime('%Y-%m-%dT%H:%M:%fZ','now') WHERE id = ? RETURNING *;

UPDATE videos SET status = ?, error_msg = ?, updated_at = strftime('%Y-%m-%dT%H:%M:%fZ','now') WHERE id = ? RETURNING *;

DELETE FROM videos WHERE id = ?;




-- queries for transcripts table  

INSERT INTO transcripts (id, video_id, language, model, status) VALUES (?, ?, ?, 'small', 'pending') RETURNING *;

-- selecting a single transcript 

SELECT * FROM transcripts WHERE id = ? LIMIT 1;

-- selecting single transcript using video id  

SELECT * FROM transcripts WHERE video_id = ? LIMIT 1;


-- update 

UPDATE transcripts SET full_text = ?, transcript_path = ?, status = 'done',  updated_at = strftime('%Y-%m-%dT%H:%M:%fZ','now') WHERE id = ? RETURNING *;

-- update status 

UPDATE transcripts SET status = ?, error_msg = ?, updated_at = strftime('%Y-%m-%dT%H:%M:%fZ','now') WHERE id = ? RETURNING *;




-- segments queries 

INSERT INTO segments (transcript_id, start_time, end_time, text, confidence) VALUES (?, ?, ?, ?, ?) RETURNING *;


-- with transcript id , ordering from their start_time (asc)

SELECT * FROM segments WHERE transcript_id = ? ORDER BY start_time ASC;


-- delete segments  

DELETE FROM segments WHERE transcript_id = ?;


-- clips queries 

INSERT INTO clips (id, video_id, clip_path, start_time, end_time, label, status) VALUES (?, ?, ?, ?, ?, ?, 'pending') RETURNING *;

SELECT * FROM clips WHERE id = ? LIMIT 1;

SELECT * FROM clips WHERE video_id = ? ORDER BY start_time ASC;

UPDATE clips SET status = ?, error_msg = ?, updated_at = strftime('%Y-%m-%dT%H:%M:%fZ', 'now') WHERE id = ? RETURNING *;

DELETE FROM clips WHERE id = ?;


-- captions queries 

INSERT INTO captions (id, clip_id, format, caption_path, burned_id) VALUES (?, ?, 'srt', ?, ?) RETURNING *;

SELECT * FROM captions WHERE id = ? LIMIT 1;

SELECT * FROM captions WHERE clip_id = ? ORDER BY created_at ASC;

DELETE FROM captions WHERE id = ?;




-- query for jobs  

-- enqueue the job

INSERT INTO jobs (id, job_type, payload, priority) VALUES (?, ?, ?, ?) RETURNING *;

SELECT * FROM jobs WHERE id = ? LIMIT 1;


-- dequeue the job - update

UPDATE jobs SET status = 'running', started_at = strftime('%Y-%m-%dT%H:%M:%fZ','now'), attempts = attempts + 1 WHERE id = (
   SELECT id FROM jobs WHERE status = 'queued' AND attempts < max_attempts ORDER BY priority DESC, queued_at ASC LIMIT 1
) RETURNING *;


-- completed job query  

UPDATE jobs SET status = 'done', ended_at = strftime('%Y-%m-%dT%H:%M:%fZ','now') WHERE id = ? RETURNING *;


-- failed job 

UPDATE jobs SET status = 'failed', last_error = ?, ended_at = strftime('%Y-%m-%dT%H:%M:%fZ','now') WHERE id = ? RETURNING *;


-- for retry - again enqueue  the job  

UPDATE jobs SET status = 'queued' WHERE id = ?, RETURNING *;

SELECT * FROM jobs WHERE status = ? ORDER BY queued_at DESC;

-- cancelled jobs 

UPDATE jobs SET status = 'cancelled', ended_at = strftime('%Y-%m-%dT%H:%M:%fZ','now') WHERE id = ? RETURNING *;


