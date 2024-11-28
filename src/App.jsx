import Layout from './components/Layout';

function App() {
  return (
    <Router>
      <Layout>
        <Routes>
          <Route path="/login" element={<Login />} />
          {/* Other routes */}
        </Routes>
      </Layout>
    </Router>
  );
}

export default App; 