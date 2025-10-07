import { ReactNode } from "react";
import styled from "styled-components";

const Wrapper = styled.div`
  background-color: rgba(0, 0, 0, 0.3);
  display: flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  height: 100%;
  position: absolute;
  top: 0;
  left: 0;
`;
const Modal = styled.div`
  display: flex;
  flex-direction: column;
  justify-content: space-around;
  gap: 0.8em;
  background-color: #777777;
  border-radius: 1em;
  padding: 1.5em;
`;
const ModalTitle = styled.h2`
  margin: 0;
`;

type PopinProps = {
  title?: string;
  isOpen: boolean;
  onClose: () => void;
  children?: ReactNode;
};
export function Popin({ children, title, isOpen, onClose }: PopinProps) {
  return isOpen ? (
    <Wrapper onClick={onClose}>
      <Modal onClick={(e) => e.stopPropagation()}>
        {!!title && <ModalTitle>{title}</ModalTitle>}
        {children}
      </Modal>
    </Wrapper>
  ) : null;
}
