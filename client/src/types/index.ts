import { ReactElement } from 'react';

export type User = {
  id: string;
  username: string;
  discriminator: string;
  global_name?: string;
  avatar?: string
}

export type Guild = {
  id: string;
  name: string;
  icon: string;
}

export type Channel = {
  id: string;
  name?: string;
  type: number;
  last_message_id?: string;
  message_count?: number;
  guild_id?: string;
  parent_id?: string;
}

export type Action = {
  type: 'CRAWL' | 'DELETE';
  scope: 'CHANNEL' | 'GUILD' | 'ALL';
  targetId?: string;
}

export type InfoListFieldConfig<T extends User | Guild | Channel> = {
  fieldName: keyof T;
  label: string;
  display?: (value?: string | number, isId?: boolean) => string | ReactElement
}

export enum LogLevel {
  DEBUG = 'DEBUG',
  INFO = 'INFO',
  SUCCESS = 'SUCCESS',
  WARNING = 'WARNING',
  ERROR = 'ERROR',
  FATAL = 'FATAL',
}

export type LogEntry = {
  message: string;
  logLevel?: LogLevel | null;
};
