import Image from "next/image";

export const dynamic = "force-dynamic";

type PostType = {
  id: number;
  title: string;
  image: string;
  content: string;
  created_at: string; // ISO Timestamp
  updated_at: string; // ISO Timestamp
};

type ApiResult<T> = {
  data: T;
  meta: {
    hostname: string;
    timestamp: string; // ISO Timestamp
  };
};

export default async function PostsPage() {
  const data = (await (
    await fetch(`${process.env.BACKEND_URL}/api/posts`)
  ).json()) as ApiResult<PostType[]>;

  return (
    <div className="flex justify-center">
      <main className="max-w-[400px] p-4">
        <div>
          {data.data.map((item) => (
            <div key={item.id} className="my-16">
              {!!item.image && (
                <Image
                  className="max-w-full"
                  src={item.image}
                  width={800}
                  height={600}
                  alt=""
                />
              )}
              <h2 className="text-2xl">
                {item.id}: {item.title}
              </h2>
              <div>
                <small>{item.created_at}</small>
              </div>
              <div>{item.content}</div>
            </div>
          ))}
        </div>
      </main>
    </div>
  );
}
