export interface UnknownMessage {
  type: string;
  data: any;
}
export interface ChatMessage {
  sender: string;
  message: string;
  owner: boolean;
}
export interface InitMessage {
  video_id: string;
  host: boolean;
  paused: boolean;
}

export interface VideoUpdateMessage {
  paused: boolean;
  time: number;
}
