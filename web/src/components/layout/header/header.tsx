import Logo from "@/assets/images/logo/logo.svg";
import ConnectButton from "@/components/common/wallet/connectButton/connectButton";
import { Flex } from "antd";
import Image from "next/image";
import Link from "next/link";

const Header = () => {
	return (
		<Flex
			justify="space-between"
			style={{
				padding: "16px 24px",
			}}>
			<Link
				href={"/"}
				style={{
					lineHeight: 0,
					textDecoration: "none",
				}}>
				<Image
					src={Logo}
					alt="Logo"
					width={93}
					height={24}
				/>
			</Link>
			<ConnectButton />
		</Flex>
	);
};

export default Header;
