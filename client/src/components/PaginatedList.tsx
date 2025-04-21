import { useEffect, useMemo, useState } from 'react';
import styled from 'styled-components';

import { Button } from '.';
import { Channel, Guild } from '../types';

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
  text-align: left;
`;

const SectionTitle = styled.h3`
  margin: 0;
`;

const Header = styled.div`
  display: flex;
  justify-content: space-between;
`;

const PAGE_SIZE = 20;

type PaginatedListProps = {
  resultsList?: Array<Guild> | Array<Channel> | null;
};

export function PaginatedList({ resultsList }: PaginatedListProps) {
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
          <li>{el.name}</li>
        ))}
      </Ul>
    </Wrapper>
  );
}
