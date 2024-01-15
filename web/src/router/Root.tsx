import { useEffect, useLayoutEffect } from 'react';
import { Outlet, useLocation } from 'react-router-dom';
// import { useStore } from 'zustand';
// import Alarm from '@src/components/common/Alarm';
// import useHistoryStack from '@src/hooks/useHistoryStack';
// import useNative from '@src/hooks/useNative';
// import RouteTracker from '@src/router/RouteTracker';
// import { useGetUserInfo } from '@src/services/query/users';
// import userPersist from '@src/store/persist';
// import { checkForUpdates } from '@src/utils/update';

const Root = () => {
	const location = useLocation();
	// 프로젝트 업데이트 체크
	useLayoutEffect(() => {
		// checkForUpdates();
	}, [location]);

	// // 히스토리 스택
	// useHistoryStack();

	//   // Native 웹뷰로부터 토큰을 받아옴
	//   {
	//     const { responseToken } = useNative();

	//     useEffect(() => {
	//       // JavaScript 함수 호출이 완료되면 해당 함수가 실행됨
	//       window.responseToken = responseToken;
	//       return () => {
	//         // 컴포넌트 언마운트 시에 함수 정리(clean-up) 등을 수행할 수 있음
	//         window.responseToken = null;
	//       };
	//     }, [responseToken]);
	//   }

	//   // Google Analytics
	//   RouteTracker();

	//   const { userId, isLogin, setLogout, setIsAccountTrue, setIsMobileTrue } =
	//     useStore(userPersist);
	//   const { data } = useGetUserInfo(userId);

	useEffect(() => {
		localStorage.setItem('currentPath', location.pathname);
	}, [location]);

	//   useEffect(() => {
	//     // 유저 아이디가 다를때 로그아웃
	//     if (userId && data?.userId) {
	//       if (userId !== data.userId) {
	//         setLogout();
	//       }
	//     }
	//   }, [userId, data, setLogout, isLogin, location]);

	//   // 유저 정보가 있을때 인증 여부 확인 후 store 업데이트
	//   useEffect(() => {
	//     if (userId && data?.bankCheck) {
	//       setIsAccountTrue();
	//     }
	//     if (userId && data?.phoneCheck) {
	//       setIsMobileTrue();
	//     }
	//   }, [userId, data, setIsAccountTrue, setIsMobileTrue, location]);

	//   useEffect(() => {
	//     if ('serviceWorker' in navigator) {
	//       navigator.serviceWorker.ready.then((registration) => {
	//         registration.unregister().then(() => {
	//           window.location.reload();
	//         });
	//       });
	//       navigator.serviceWorker.getRegistrations().then((registrations) => {
	//         registrations.forEach((registration) => registration.unregister());
	//       });
	//     }
	//     if ('caches' in window) {
	//       caches.keys().then((keyList) => {
	//         return Promise.all(
	//           keyList.map((key) => {
	//             return caches.delete(key);
	//           }),
	//         );
	//       });
	//     }
	//   }, []);

	return (
		<>
			<Outlet />
			{/* <Alarm /> */}
		</>
	);
};

export default Root;
