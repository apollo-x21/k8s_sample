import { useEffect, useState } from 'react';
import type { AuthPayload } from '../services/backend';
import { getUsers, login, logout, register } from '../services/backend';
import './index.less';

type Mode = 'login' | 'register';

export default function IndexPage() {
  const [mode, setMode] = useState<Mode>('login');
  const [form, setForm] = useState<AuthPayload>({ username: '', password: '' });
  const [token, setToken] = useState<string>();
  const [users, setUsers] = useState<string[]>([]);
  const [currentUser, setCurrentUser] = useState('');
  const [status, setStatus] = useState('');
  const [error, setError] = useState('');

  useEffect(() => {
    if (!token) {
      return;
    }
    getUsers(token)
      .then((data) => {
        setUsers(data.users.map((u) => u.username));
        setCurrentUser(data.me);
        setError('');
        setStatus('用户列表已更新');
      })
      .catch((err) => setError(err.message));
  }, [token]);

  const updateField = (key: keyof AuthPayload, value: string) => {
    setForm((prev) => ({ ...prev, [key]: value }));
  };

  const handleSubmit = async (event: React.FormEvent) => {
    event.preventDefault();
    setStatus('');
    setError('');

    try {
      if (mode === 'register') {
        const result = await register(form);
        setStatus(result.message);
        setMode('login');
      } else {
        const result = await login(form);
        if (result.token) {
          setToken(result.token);
          setStatus('登录成功，正在获取用户列表');
        }
      }
    } catch (err) {
      if (err instanceof Error) {
        setError(err.message);
      } else {
        setError('未知错误');
      }
    }
  };

  const authenticated = Boolean(token);

  const handleLogout = async () => {
    if (!token) {
      return;
    }
    setError('');
    setStatus('');
    try {
      const response = await logout(token);
      setStatus(response.message || '已退出登录');
    } catch (err) {
      if (err instanceof Error) {
        setError(err.message);
        return;
      }
      setError('退出失败');
      return;
    }
    setToken(undefined);
    setUsers([]);
    setCurrentUser('');
  };

  const rootClass = authenticated ? 'app admin' : 'app auth';

  return (
    <div className={rootClass}>
      {!authenticated ? (
        <div className="auth-card">
          <h1>微服务示例前端</h1>
          <div className="toggle">
            <button
              type="button"
              className={mode === 'login' ? 'active' : ''}
              onClick={() => setMode('login')}
            >
              登录
            </button>
            <button
              type="button"
              className={mode === 'register' ? 'active' : ''}
              onClick={() => setMode('register')}
            >
              注册
            </button>
          </div>
          <form onSubmit={handleSubmit} className="form">
            <label>
              用户名
              <input
                value={form.username}
                onChange={(e) => updateField('username', e.target.value)}
                placeholder="输入用户名"
              />
            </label>
            <label>
              密码
              <input
                type="password"
                value={form.password}
                onChange={(e) => updateField('password', e.target.value)}
                placeholder="输入密码"
              />
            </label>
            <button className="primary" type="submit">
              {mode === 'login' ? '登录' : '注册'}
            </button>
          </form>
          {status && <div className="status">{status}</div>}
          {error && <div className="error">{error}</div>}
        </div>
      ) : (
        <div className="layout">
          <aside className="sidebar">
            <div className="logo">微服务控制台</div>
            <nav className="nav">
              <button className="nav-item active">用户列表</button>
            </nav>
          </aside>
          <section className="content">
            <header className="topbar">
              <div className="title">用户中心</div>
              <button className="profile" type="button" onClick={handleLogout}>
                {currentUser} · 退出
              </button>
            </header>
            <main className="main">
              <div className="card">
                <div className="me">当前用户：{currentUser}</div>
                <h2>系统用户</h2>
                <table className="user-table">
                  <thead>
                    <tr>
                      <th>#</th>
                      <th>用户名</th>
                    </tr>
                  </thead>
                  <tbody>
                    {users.map((user, index) => (
                      <tr key={user}>
                        <td>{index + 1}</td>
                        <td>{user}</td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
              <div className="feedback">
                {status && <div className="status">{status}</div>}
                {error && <div className="error">{error}</div>}
              </div>
            </main>
          </section>
        </div>
      )}
    </div>
  );
}
