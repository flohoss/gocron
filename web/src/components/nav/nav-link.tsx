import type { PropFunction } from '@builder.io/qwik';
import { component$ } from '@builder.io/qwik';
import { Link } from '@builder.io/qwik-city';

export default component$(
  (props: {
    link: string;
    name: string;
    extra?: string;
    icon: string;
    onClick$?: PropFunction<() => void>;
    status?: string;
    active: boolean;
    bg?: 'bg-primary' | 'bg-secondary';
    text?: 'text-primary' | 'text-secondary';
  }) => {
    return (
      <li>
        <Link onClick$={props.onClick$} href={props.link} class={`flex justify-between px-2 ${props.active && 'active'}`}>
          <div class="flex items-center gap-2">
            <div
              class={`w-10 h-10 p-2 rounded-full ${props.bg ? props.bg : 'bg-base-100'} flex items-center justify-center shrink-0 text-lg ${
                props.text ? props.text : 'text-base'
              }`}
              dangerouslySetInnerHTML={props.icon}
            ></div>
            <div class="flex flex-col">
              <div class="text-lg">{props.name}</div>
              {props.extra && <div class="opacity-50 text-xs">{props.extra}</div>}
            </div>
          </div>
          {props.status && (
            <div class="relative flex h-3 w-3">
              <span class="animate-ping absolute inline-flex h-full w-full rounded-full bg-success opacity-75"></span>
              <span class="relative inline-flex rounded-full h-3 w-3 bg-success"></span>
            </div>
          )}
        </Link>
      </li>
    );
  }
);
