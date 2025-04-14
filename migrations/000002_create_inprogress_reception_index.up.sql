CREATE UNIQUE INDEX idx_unique_inprogress_reception
ON reception (pvz_id)
WHERE status = 'in_progress';