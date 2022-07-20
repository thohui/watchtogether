export interface UnknownMessage {
  type: string;
  data: any;
}

export interface ChatMessage {
  sender: string;
  message: string;
}
