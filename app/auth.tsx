import { createContext, useContext, useEffect, useState } from 'react';
import { initializeApp, getApps, getApp } from 'firebase/app';
import { getFirestore, doc, setDoc } from 'firebase/firestore/lite';
import {
  getAuth,
  onAuthStateChanged,
  GithubAuthProvider,
  signInWithPopup,
  signOut,
  User,
  OAuthCredential,
} from 'firebase/auth';

// Your web app's Firebase configuration
// For Firebase JS SDK v7.20.0 and later, measurementId is optional
const firebaseConfig = {
  apiKey: process.env.NEXT_PUBLIC_FIREBASE_API_KEY,
  authDomain: process.env.NEXT_PUBLIC_FIREBASE_AUTH_DOMAIN,
  databaseURL: process.env.NEXT_PUBLIC_FIREBASE_DATABASE_URL,
  projectId: process.env.NEXT_PUBLIC_FIREBASE_PROJECT_ID,
  storageBucket: process.env.NEXT_PUBLIC_FIREBASE_STORAGE_BUCKET,
  // messagingSenderId: process.env.NEXT_PUBLIC_FIREBASE_MESSAGING_SENDER_ID,
  appId: process.env.NEXT_PUBLIC_FIREBASE_APP_ID,
  // measurementId: process.env.NEXT_PUBLIC_FIREBASE_MEASUREMENT_ID
};
const app = getApps().length === 0 ? initializeApp(firebaseConfig) : getApp();
export const auth = getAuth(app);
export const firestore = getFirestore(app);
const ctx = createContext<ReturnType<typeof useFirebaseAuth> | null>(null);

export function AuthProvider({ children }: { children: React.ReactNode }) {
  const firebaseAuth = useFirebaseAuth();
  return <ctx.Provider value={firebaseAuth}>{children}</ctx.Provider>;
}

export function useAuth() {
  return useContext(ctx);
}

function useFirebaseAuth() {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState<boolean>(true);

  function handleUser(
    user: User | null,
    credential: OAuthCredential | null = null,
  ) {
    if (user) {
      setLoading(false);
      setDoc(
        doc(firestore, 'users', user.uid),
        {
          displayName: user.displayName,
          email: user.email,
          phoneNumber: user.phoneNumber,
          photoURL: user.photoURL,
          uid: user.uid,
          accessToken: credential?.accessToken || null,
        },
        { mergeFields: ['displayName', 'email', 'phoneNumber', 'photoURL'] },
      );
      return user;
    } else {
      setUser(null);
      setLoading(false);
      return null;
    }
  }

  function signinWithGithub() {
    setLoading(true);
    return signInWithPopup(
      auth,
      new GithubAuthProvider().addScope('public_repo'),
    ).then((response) =>
      handleUser(
        response.user,
        GithubAuthProvider.credentialFromResult(response),
      ),
    );
  }

  function signout() {
    return signOut(auth).then(() => handleUser(null));
  }

  useEffect(() => {
    const unsubscribe = onAuthStateChanged(auth, handleUser);
    return () => unsubscribe();
  }, []);

  return { user, loading, signinWithGithub, signout };
}
