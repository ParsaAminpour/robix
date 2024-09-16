import WalletModal from "@/components/common/wallet/walletModal/walletModal";
import { triggerModal } from "@/store/slices/modal/modal.slice";
import { useDispatch, useSelector } from "@/store/store";
import React from "react";

const ModalProvider = () => {
	const dispatch = useDispatch();
	const { modals } = useSelector((state) => state.modal);

	return (
		<>
			<WalletModal
				onClose={() => dispatch(triggerModal({ modal: "wallet", trigger: false }))}
				open={modals.wallet}
			/>
		</>
	);
};

export default ModalProvider;
