import { Breadcrumb } from 'react-bootstrap';
import { useLocation, Link } from 'react-router-dom';

const Breadcrumbs = () => {
  const location = useLocation();
  const pathnames = location.pathname.split('/').filter(x => x);

  return (
    <Breadcrumb className="mb-4">
      <Breadcrumb.Item linkAs={Link} linkProps={{ to: '/' }}>Home</Breadcrumb.Item>
      {pathnames.map((name, index) => {
        const routeTo = `/${pathnames.slice(0, index + 1).join('/')}`;
        const displayName = name.charAt(0).toUpperCase() + name.slice(1).replace('-', ' ');
        return (
          <Breadcrumb.Item key={routeTo} linkAs={Link} linkProps={{ to: routeTo }}>
            {displayName}
          </Breadcrumb.Item>
        );
      })}
    </Breadcrumb>
  );
};

export default Breadcrumbs;