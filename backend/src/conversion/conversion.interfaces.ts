export interface ResponseMessage {
  operation_id: string;
  status: 'progress' | 'error' | 'success';
  message?: string;
  progress?: number;
  new_file_path?: string;
}
