import { Link } from "react-router-dom";

export function NotFoundPage() {
  return (
    <section>
      <h1>見つかりません</h1>
      <p>
        お探しのページは存在しません。<Link to="/">一覧へ戻る</Link>
      </p>
    </section>
  );
}
