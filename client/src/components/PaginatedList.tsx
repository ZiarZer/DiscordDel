import { useCallback, useEffect, useMemo, useState } from 'react';
import styled, { css } from 'styled-components';

import { Button } from '.';
import { Channel, Guild } from '../types';
import CopyIcon from '../assets/copy.svg?react';
import ExternalLinkIcon from '../assets/external-link.svg?react';

const Wrapper = styled.div`
  background-color: #ffffff30;
  border-radius: 1em;
  padding: 1em;
  text-align: center;
  display: flex;
  flex-direction: column;
  gap: 0.25em;
  width: 33%;
`;

const PageNav = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-evenly;
  flex: 1;
  gap: 0.25em;
`;

const Ul = styled.ul`
  list-style-type: none;
  text-align: left;
  padding: unset;
  li {
    display: flex;
    align-items: center;
    gap: 0.25em;
    margin: 0.25em auto;
  }
`;

const SectionTitle = styled.h3`
  margin: 0;
`;

const Header = styled.div`
  display: flex;
  justify-content: space-between;
`;

const iconContainerStyle = css`
  display: flex;
  border-radius: 4px;
  align-items: center;
  background-color: #1a1a1a;
  padding: 0.25em;
  color: inherit;
  svg {
    height: 1em;
    width: 1em;
  }
`;

const IconButton = styled.div`
  ${iconContainerStyle}
  cursor: copy;
`;

const IconLink = styled.a`
  ${iconContainerStyle}
`;

const PAGE_SIZE = 20;

type PaginatedListProps = {
  resultsList?: Array<Guild> | Array<Channel> | null;
  isChannelType?: boolean;
};

export function PaginatedList({
  resultsList,
  isChannelType = false,
}: PaginatedListProps) {
  const [page, setPage] = useState(0);
  const pageCount = useMemo(
    () =>
      resultsList != null ? Math.ceil(resultsList.length / PAGE_SIZE) : null,
    [resultsList]
  );
  const hasPrevious = useMemo(() => page > 0, [page]);
  const hasNext = useMemo(
    () => pageCount && page + 1 < pageCount,
    [page, pageCount]
  );
  const nextPage = () => {
    if (hasNext) {
      setPage(page + 1);
    }
  };
  const previousPage = () => {
    if (hasPrevious) {
      setPage(page - 1);
    }
  };

  const displayedResults = useMemo(
    () => resultsList?.slice(PAGE_SIZE * page, PAGE_SIZE * (page + 1)),
    [resultsList, page]
  );

  useEffect(() => setPage(0), [resultsList]);

  const getItemLink = useCallback(
    (object) => {
      if (!isChannelType) {
        return undefined;
      }
      const guildPart = object.guild_id ?? '@me';
      return `https://discord.com/channels/${guildPart}/${object.id}`;
    },
    [isChannelType]
  );

  return (
    <Wrapper>
      <SectionTitle>Results</SectionTitle>
      <Header>
        <PageNav>
          <Button disabled={!hasPrevious} onClick={previousPage}>
            Previous
          </Button>
          {displayedResults ? page + 1 : '_'} / {pageCount ?? '_'}
          <Button disabled={!hasNext} onClick={nextPage}>
            Next
          </Button>
        </PageNav>
      </Header>
      <Ul>
        {displayedResults?.map((el: Guild | Channel) => (
          <li>
            <IconButton onClick={() => navigator.clipboard.writeText(el.id)}>
              <CopyIcon />
            </IconButton>
            {isChannelType && (
              <IconLink href={getItemLink(el)} rel="noreferrer" target="_blank">
                <ExternalLinkIcon />
              </IconLink>
            )}
            {el.name}
          </li>
        ))}
      </Ul>
    </Wrapper>
  );
}
