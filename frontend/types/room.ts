export interface UnknownMessage {
  type: string;
  data: any;
}
export interface ChatMessage {
  sender: string;
  message: string;
}
export interface InitMessage {
  video_id: string;
  time: number;
  host: boolean;
  paused: boolean;
}

export interface VideoUpdateMessage {
  paused: boolean;
  time: number;
}
